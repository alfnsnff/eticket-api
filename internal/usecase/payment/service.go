package payment

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	enum "eticket-api/internal/common/enums"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PaymentUsecase struct {
	DB                     *gorm.DB // Assuming you have a DB field for the transaction manager
	TripayClient           *client.TripayClient
	ClaimSessionRepository ClaimSessionRepository
	BookingRepository      BookingRepository
	TicketRepository       TicketRepository
	Mailer                 mailer.Mailer
}

func NewPaymentUsecase(
	db *gorm.DB,
	tripay_client *client.TripayClient,
	claim_session_repository ClaimSessionRepository,
	booking_repository BookingRepository,
	ticket_repository TicketRepository,
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

func (c *PaymentUsecase) GetPaymentChannels(ctx context.Context) ([]*model.ReadPaymentChannelResponse, error) {
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

	session, err := c.ClaimSessionRepository.GetByUUIDWithLock(tx, sessionID, true)
	if err != nil {
		return nil, fmt.Errorf("get claim session failed: %w", err)
	}

	if session == nil || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("claim session not found or expired")
	}

	if session.Status != enum.ClaimSessionPendingPayment.String() {
		return nil, fmt.Errorf("claim session has invalid status: %s", session.Status)
	}

	var Amount float32
	booking, err := c.BookingRepository.GetByOrderID(tx, request.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	tickets, err := c.TicketRepository.GetByBookingID(tx, booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tickets: %w", err)
	}
	if len(tickets) == 0 {
		return nil, fmt.Errorf("no tickets found for booking %s", *booking.OrderID)
	}

	// Sum ticket prices
	for _, ticket := range tickets {
		Amount += ticket.Price
	}

	// Create payment
	response, err := c.TripayClient.CreatePayment(
		request.PaymentMethod,
		int(Amount),
		booking.CustomerName,
		booking.Email,
		booking.PhoneNumber,
		*booking.OrderID,
		tickets,
	)
	if err != nil {
		return nil, fmt.Errorf("create Tripay payment failed: %w", err)
	}

	// Update booking with reference number
	if response.Reference != "" {
		ref := response.Reference // make a pointer
		booking.ReferenceNumber = &ref
		if err := c.BookingRepository.UpdateReferenceNumber(tx, booking.ID, &ref); err != nil {
			return nil, fmt.Errorf("failed to update booking with reference number: %w", err)
		}
	}
	session = &entity.ClaimSession{
		ID:         session.ID,
		SessionID:  session.SessionID,
		ScheduleID: session.ScheduleID,
		Status:     enum.ClaimSessionPendingTransaction.String(),
		ExpiresAt:  time.Unix(int64(response.ExpiredTime), 0), // Convert Unix timestamp to time.Time
	}

	if err := c.ClaimSessionRepository.Update(tx, session); err != nil {
		return nil, fmt.Errorf("update claim session failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &response, nil
}

func (c *PaymentUsecase) HandleCallback(ctx context.Context, r *http.Request, request *model.WriteCallbackRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	booking, err := c.BookingRepository.GetByOrderID(tx, request.MerchantRef)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	tickets, err := c.TicketRepository.GetByBookingID(tx, booking.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve tickets: %w", err)
	}
	if len(tickets) == 0 { // Changed from tickets == nil to len(tickets) == 0
		return fmt.Errorf("no tickets found for booking")
	}

	// Get session from the first ticket's claim_session_id
	var session *entity.ClaimSession
	if tickets[0].ClaimSessionID != nil {
		session, err = c.ClaimSessionRepository.GetByID(tx, *tickets[0].ClaimSessionID)
		if err != nil {
			return fmt.Errorf("failed to get claim session: %w", err)
		}
		if session == nil {
			return fmt.Errorf("claim session not found")
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
func (c *PaymentUsecase) HandleSuccessfulPayment(tx *gorm.DB, booking *entity.Booking, tickets []*entity.Ticket, session *entity.ClaimSession) error {
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
func (c *PaymentUsecase) HandleUnsuccessfulPayment(tx *gorm.DB, booking *entity.Booking, tickets []*entity.Ticket, session *entity.ClaimSession, status string) error {

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
