package usecase

import (
	"context"
	enum "eticket-api/internal/common/enums"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PaymentUsecase struct {
	DB                     *gorm.DB // Assuming you have a DB field for the transaction manager
	TripayClient           domain.TripayClient
	ClaimSessionRepository domain.ClaimSessionRepository
	BookingRepository      domain.BookingRepository
	TicketRepository       domain.TicketRepository
	Mailer                 mailer.Mailer
}

func NewPaymentUsecase(
	db *gorm.DB,
	tripay_client domain.TripayClient,
	claim_session_repository domain.ClaimSessionRepository,
	booking_repository domain.BookingRepository,
	ticket_repository domain.TicketRepository,
	mailer mailer.Mailer,
) *PaymentUsecase {
	return &PaymentUsecase{
		DB:                     db,
		TripayClient:           tripay_client,
		ClaimSessionRepository: claim_session_repository,
		BookingRepository:      booking_repository,
		TicketRepository:       ticket_repository,
		Mailer:                 mailer,
	}
}

func (c *PaymentUsecase) ListPaymentChannels(ctx context.Context) ([]*model.ReadPaymentChannelResponse, error) {
	channels, err := c.TripayClient.GetPaymentChannels()
	if err != nil {
		return nil, err
	}
	result := make([]*model.ReadPaymentChannelResponse, len(channels))
	for i := range channels {
		result[i] = &channels[i]
	}
	return result, nil
}

func (c *PaymentUsecase) GetTransactionDetail(ctx context.Context, reference string) (*model.ReadTransactionResponse, error) {
	return c.TripayClient.GetTransactionDetail(reference)
}

func (c *PaymentUsecase) CreatePayment(ctx context.Context, request *model.WritePaymentRequest, sessionID string) (*model.ReadTransactionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	session, err := c.ClaimSessionRepository.FindBySessionID(tx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if session == nil {
		return nil, errs.ErrNotFound
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errs.ErrExpired
	}

	booking, err := c.BookingRepository.FindByOrderID(tx, request.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errs.ErrNotFound
	}

	tickets, err := c.TicketRepository.FindByBookingID(tx, booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets: %w", err)
	}
	if len(tickets) == 0 {
		return nil, errs.ErrNotFound
	}

	var amounts float64
	for _, ticket := range tickets {
		amounts += ticket.Price
	}

	orderItems := make([]model.OrderItem, len(tickets))
	for i, ticket := range tickets {
		orderItems[i] = mapper.TicketToItem(ticket)
	}

	payload := &model.WriteTransactionRequest{
		Method:        request.PaymentMethod,
		Amount:        int(amounts), // Convert to integer cents
		CustomerName:  booking.CustomerName,
		CustomerEmail: booking.Email,
		CustomerPhone: booking.PhoneNumber,
		MerchantRef:   *booking.OrderID,
		OrderItems:    orderItems,
		CallbackUrl:   "https://example.com/callback",
		ReturnUrl:     "https://example.com/callback",
		ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),
	}

	// Create payment
	payment, err := c.TripayClient.CreatePayment(payload)
	if err != nil {
		return nil, fmt.Errorf("create Tripay payment failed: %w", err)
	}

	booking.ReferenceNumber = &payment.Reference
	if err := c.BookingRepository.Update(tx, booking); err != nil {
		return nil, fmt.Errorf("failed to update booking with reference number: %w", err)
	}
	session = &domain.ClaimSession{
		ID:         session.ID,
		SessionID:  session.SessionID,
		ScheduleID: session.ScheduleID,
		Status:     enum.ClaimSessionPendingTransaction.String(),
		ExpiresAt:  time.Unix(int64(payment.ExpiredTime), 0), // Convert Unix timestamp to time.Time
	}

	if err := c.ClaimSessionRepository.Update(tx, session); err != nil {
		return nil, fmt.Errorf("update claim session failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &payment, nil
}

func (c *PaymentUsecase) HandleCallback(ctx context.Context, request *model.WriteCallbackRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	booking, err := c.BookingRepository.FindByOrderID(tx, request.MerchantRef)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return errs.ErrNotFound
	}
	tickets, err := c.TicketRepository.FindByBookingID(tx, booking.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve tickets: %w", err)
	}
	if len(tickets) == 0 { // Changed from tickets == nil to len(tickets) == 0
		return errs.ErrNotFound
	}

	// Get session from the first ticket's claim_session_id
	var session *domain.ClaimSession
	if tickets[0].ClaimSessionID != nil {
		session, err = c.ClaimSessionRepository.FindByID(tx, *tickets[0].ClaimSessionID)
		if err != nil {
			return fmt.Errorf("failed to get session: %w", err)
		}
		if session == nil {
			return errs.ErrNotFound
		}
	} else {
		return fmt.Errorf("tickets have no associated claim session")
	}

	// Handle different payment statuses
	switch request.Status {
	case "PAID":
		// Payment successful
		if err := c.HandleSuccessfulPayment(tx, booking, tickets, session); err != nil {
			return err
		}

	case "FAILED", "EXPIRED", "REFUND":
		// Payment unsuccessful
		if err := c.HandleUnsuccessfulPayment(tx, booking, tickets, session, request.Status); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown payment status: %s", request.Status)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Handle successful payment
func (c *PaymentUsecase) HandleSuccessfulPayment(tx *gorm.DB, booking *domain.Booking, tickets []*domain.Ticket, session *domain.ClaimSession) error {
	// Update session to success status
	session.Status = enum.ClaimSessionSuccess.String()

	if err := c.ClaimSessionRepository.Update(tx, session); err != nil {
		return fmt.Errorf("update claim session failed: %w", err)
	}

	// Send confirmation email
	subject := "Your Booking is Confirmed"
	htmlBody := templates.BookingSuccessEmail(booking.CustomerName, *booking.OrderID, len(tickets), time.Now().Year())

	if err := c.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return nil
}

// Handle unsuccessful payment
func (c *PaymentUsecase) HandleUnsuccessfulPayment(tx *gorm.DB, booking *domain.Booking, tickets []*domain.Ticket, session *domain.ClaimSession, status string) error {

	// Update session to failed status
	session.Status = enum.ClaimSessionFailed.String() // Expire immediately

	if err := c.ClaimSessionRepository.Update(tx, session); err != nil {
		return fmt.Errorf("update claim session failed: %w", err)
	}

	subject := "Payment Failed - Booking Not Confirmed"
	var htmlBody string

	switch status {
	case "FAILED":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment processing failed")
	case "EXPIRED":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment time expired")
	case "CANCELLED":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment was cancelled")
	}

	if err := c.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
		fmt.Printf("Warning: failed to send failure email: %v", err)
	}

	return nil
}
