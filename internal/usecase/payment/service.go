package payment

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/model"
	"eticket-api/internal/repository"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PaymentUsecase struct {
	DB                *gorm.DB // Assuming you have a DB field for the transaction manager
	TripayClient      *client.TripayClient
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
	Mailer            *mailer.SMTPMailer
}

func NewPaymentUsecase(
	db *gorm.DB,
	tripay_client *client.TripayClient,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
	mailer *mailer.SMTPMailer,
) *PaymentUsecase {
	return &PaymentUsecase{
		DB:                db,
		TripayClient:      tripay_client,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
		Mailer:            mailer,
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

	if request.OrderID == "" {
		return nil, errors.New("missing required field: OrderID")
	}
	var Amount float32

	var err error
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
		return nil, fmt.Errorf("no tickets found for booking %s", booking.OrderID)
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
		booking.OrderID,
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

	if request.MerchantRef == "" {
		return errors.New("invalid request: all customer fields are required")
	}

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
	if tickets == nil {
		return fmt.Errorf("tiket is empty not found")
	}

	for _, ticket := range tickets {
		c.TicketRepository.Paid(tx, ticket.ID)
	}

	subject := "Your Booking is Confirmed"

	// Inside your callback logic:
	htmlBody := templates.BookingSuccessEmail(booking.CustomerName, booking.OrderID, len(tickets), time.Now().Year())

	if err := c.Mailer.Send(booking.Email, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
