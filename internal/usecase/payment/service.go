package payment

import (
	"context"
	"errors"
	"eticket-api/internal/client"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/repository"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type PaymentUsecase struct {
	Tx                *tx.TxManager
	TripayClient      *client.TripayClient
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
}

func NewPaymentUsecase(
	tx *tx.TxManager,
	tripay_client *client.TripayClient,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
) *PaymentUsecase {
	return &PaymentUsecase{
		Tx:                tx,
		TripayClient:      tripay_client,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
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

func (c *PaymentUsecase) CreatePayment(ctx context.Context, request *model.WritePaymentRequest) (*model.ReadTransactionResponse, error) {
	if request.OrderID == "" {
		return nil, errors.New("invalid request: all customer fields are required")
	}
	Booking := new(entity.Booking)
	Tickets := []*entity.Ticket{}
	var Amount float32
	if err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		booking, err := c.BookingRepository.GetByOrderID(tx, request.OrderID)
		if err != nil {
			return err
		}
		if booking == nil {
			return fmt.Errorf("booking not found")
		}

		Booking = booking
		tickets, err := c.TicketRepository.GetByBookingID(tx, booking.ID)

		if err != nil {
			return err
		}
		if tickets == nil {
			return fmt.Errorf("tiket is empty not found")
		}
		Tickets = tickets
		return err
	}); err != nil {
		return nil, err
	}

	for _, ticket := range Tickets {
		Amount += ticket.Price
	}

	response, err := c.TripayClient.CreatePayment(request.PaymentMethod, int(Amount), Booking.CustomerName, Booking.Email, Booking.PhoneNumber, Booking.OrderID, Tickets)
	if err != nil {
		return nil, fmt.Errorf("create Tripay payment failed: %w, tickets: %+v", err, Tickets)
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

		return err
	})

	if err != nil {
		return err
	}
	return err
}
