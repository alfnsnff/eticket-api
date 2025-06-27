package usecase

import (
	"context"
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
	DB                *gorm.DB // Assuming you have a DB field for the transaction manager
	TripayClient      domain.TripayClient
	BookingRepository domain.BookingRepository
	TicketRepository  domain.TicketRepository
	QuotaRepository   domain.QuotaRepository
	Mailer            mailer.Mailer
}

func NewPaymentUsecase(
	db *gorm.DB,
	tripay_client domain.TripayClient,
	claim_session_repository domain.ClaimSessionRepository,
	booking_repository domain.BookingRepository,
	ticket_repository domain.TicketRepository,
	quota_repository domain.QuotaRepository,
	mailer mailer.Mailer,
) *PaymentUsecase {
	return &PaymentUsecase{
		DB:                db,
		TripayClient:      tripay_client,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
		QuotaRepository:   quota_repository,
		Mailer:            mailer,
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

func (c *PaymentUsecase) TESTCreatePayment(ctx context.Context, request *model.WritePaymentRequest, orderID string) (*model.ReadTransactionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	booking, err := c.BookingRepository.FindByOrderID(tx, request.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, errs.ErrNotFound
	}

	var amounts float64
	for _, ticket := range booking.Tickets {
		amounts += ticket.Price
	}

	orderItems := make([]model.OrderItem, len(booking.Tickets))
	for i, ticket := range booking.Tickets {
		orderItems[i] = mapper.TicketToItem(&ticket)
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

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &payment, nil
}

func (c *PaymentUsecase) CreatePayment(ctx context.Context, request *model.WritePaymentRequest) (*model.ReadTransactionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

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

	// Handle different payment statuse
	switch request.Status {
	case "PAID":
		// Payment successful
		if err := c.HandleSuccessfulPayment(tx, booking, tickets); err != nil {
			return err
		}

	case "FAILED", "EXPIRED", "REFUND":
		// Payment unsuccessful
		if err := c.HandleUnsuccessfulPayment(tx, booking, tickets, request.Status); err != nil {
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
func (c *PaymentUsecase) HandleSuccessfulPayment(tx *gorm.DB, booking *domain.Booking, tickets []*domain.Ticket) error {
	// Send confirmation email
	subject := "Your Booking is Confirmed"
	htmlBody := templates.BookingSuccessEmail(booking.CustomerName, *booking.OrderID, len(tickets), time.Now().Year())

	if err := c.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	return nil

}

// Handle unsuccessful payment
func (c *PaymentUsecase) HandleUnsuccessfulPayment(tx *gorm.DB, booking *domain.Booking, tickets []*domain.Ticket, status string) error {
	// Restore quota for each ticket's class
	restored := make(map[uint]bool)
	for _, ticket := range tickets {
		if ticket == nil {
			continue
		}
		// Avoid double-restoring quota for the same class
		if restored[ticket.ClassID] {
			continue
		}
		quota, err := c.QuotaRepository.FindByScheduleIDAndClassID(tx, booking.ScheduleID, ticket.ClassID)
		if err == nil && quota != nil {
			// Count tickets for this class
			count := 0
			for _, t := range tickets {
				if t != nil && t.ClassID == ticket.ClassID {
					count++
				}
			}
			quota.Quota += count
			if err := c.QuotaRepository.Update(tx, quota); err != nil {
				fmt.Printf("Warning: failed to restore quota for class %d: %v\n", ticket.ClassID, err)
			}
			restored[ticket.ClassID] = true
		} else if err != nil {
			fmt.Printf("Warning: failed to find quota for class %d: %v\n", ticket.ClassID, err)
		}
	}

	subject := "Payment Failed - Booking Not Confirmed"
	var htmlBody string

	switch status {
	case "FAILED":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment processing failed")
	case "EXPIRED":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment time expired")
	case "CANCELLED", "REFUND":
		htmlBody = templates.BookingFailedEmail(booking.CustomerName, *booking.OrderID, "Payment was cancelled or refunded")
	}

	if err := c.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
		fmt.Printf("Warning: failed to send failure email: %v\n", err)
	}

	return nil
}
