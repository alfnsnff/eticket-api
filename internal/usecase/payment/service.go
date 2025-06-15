package payment

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/repository"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PaymentUsecase struct {
	Tx                *tx.TxManager
	TripayClient      *client.TripayClient
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
	Mailer            *mailer.SMTPMailer
}

func NewPaymentUsecase(
	tx *tx.TxManager,
	tripay_client *client.TripayClient,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
	mailer *mailer.SMTPMailer,
) *PaymentUsecase {
	return &PaymentUsecase{
		Tx:                tx,
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
	if request.OrderID == "" {
		return nil, errors.New("invalid request: all customer fields are required")
	}

	var Booking *entity.Booking
	var Tickets []*entity.Ticket
	var Amount float32

	// Retrieve booking and tickets
	if err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		Booking, err = c.BookingRepository.GetByOrderID(tx, request.OrderID)
		if err != nil {
			return err
		}
		if Booking == nil {
			return fmt.Errorf("booking not found")
		}

		Tickets, err = c.TicketRepository.GetByBookingID(tx, Booking.ID)
		if err != nil {
			return err
		}
		if len(Tickets) == 0 {
			return fmt.Errorf("tickets not found")
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Sum ticket prices
	for _, ticket := range Tickets {
		Amount += ticket.Price
	}

	// Create payment
	response, err := c.TripayClient.CreatePayment(
		request.PaymentMethod,
		int(Amount),
		Booking.CustomerName,
		Booking.Email,
		Booking.PhoneNumber,
		Booking.OrderID,
		Tickets,
	)
	if err != nil {
		return nil, fmt.Errorf("create Tripay payment failed: %w", err)
	}

	// Update booking with reference number
	if response.Reference != "" {
		ref := response.Reference // make a pointer
		if err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
			Booking.ReferenceNumber = &ref
			return c.BookingRepository.UpdateReferenceNumber(tx, Booking.ID, &ref)
		}); err != nil {
			return nil, fmt.Errorf("failed to update booking with reference number: %w", err)
		}
	}

	return &response, nil
}

func (c *PaymentUsecase) HandleCallback(ctx context.Context, r *http.Request, request *model.WriteCallbackRequest) error {
	if request.MerchantRef == "" {
		return errors.New("invalid request: all customer fields are required")
	}
	err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		booking, err := c.BookingRepository.GetByOrderID(tx, request.MerchantRef)
		if err != nil {
			return err
		}
		if booking == nil {
			return fmt.Errorf("booking not found")
		}

		tickets, err := c.TicketRepository.GetByBookingID(tx, booking.ID)

		if err != nil {
			return err
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

		return err
	})

	if err != nil {
		return err
	}
	return err
}
