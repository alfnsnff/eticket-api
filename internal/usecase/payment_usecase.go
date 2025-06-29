package usecase

import (
	"context"
	"eticket-api/internal/client"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"eticket-api/pkg/gotann"
	"fmt"
	"strings"
	"time"
)

type PaymentUsecase struct {
	Transactor        *transact.Transactor // Assuming transact package is imported
	TripayClient      domain.TripayClient
	BookingRepository domain.BookingRepository
	TicketRepository  domain.TicketRepository
	QuotaRepository   domain.QuotaRepository
	Mailer            mailer.Mailer
}

func NewPaymentUsecase(
	transactor *transact.Transactor, // Assuming transact package is imported
	tripay_client domain.TripayClient,
	booking_repository domain.BookingRepository,
	ticket_repository domain.TicketRepository,
	quota_repository domain.QuotaRepository,
	mailer mailer.Mailer,
) *PaymentUsecase {
	return &PaymentUsecase{
		Transactor:        transactor,
		TripayClient:      tripay_client,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
		QuotaRepository:   quota_repository,
		Mailer:            mailer,
	}
}

func (uc *PaymentUsecase) ListPaymentChannels(ctx context.Context) ([]*domain.PaymentChannel, error) {
	channels, err := uc.TripayClient.GetPaymentChannels()
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (uc *PaymentUsecase) GetTransactionDetail(ctx context.Context, reference string) (*domain.Transaction, error) {
	return uc.TripayClient.GetTransactionDetail(reference)
}

func (uc *PaymentUsecase) CreatePayment(ctx context.Context, request *model.WritePaymentRequest) (*domain.Transaction, error) {
	var err error
	var payment *domain.Transaction
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByOrderID(ctx, tx, request.OrderID)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		tickets, err := uc.TicketRepository.FindByBookingID(ctx, tx, booking.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}
		if len(tickets) == 0 {
			return errs.ErrNotFound
		}

		var amounts float64
		for _, ticket := range tickets {
			amounts += ticket.Price
		}

		orderItems := make([]domain.OrderItem, len(tickets))
		for i, ticket := range tickets {
			orderItems[i] = client.TicketToItem(ticket)
		}

		payload := &domain.TransactionRequest{
			Method:        request.PaymentMethod,
			Amount:        int(amounts), // Convert to integer cents
			CustomerName:  booking.CustomerName,
			CustomerEmail: booking.Email,
			CustomerPhone: booking.PhoneNumber,
			MerchantRef:   booking.OrderID,
			OrderItems:    orderItems,
			CallbackUrl:   "https://example.com/callback",
			ReturnUrl:     "https://example.com/callback",
			ExpiredTime:   int(time.Now().Add(30 * time.Minute).Unix()),
		}

		payment, err = uc.TripayClient.CreatePayment(payload)
		if err != nil {
			// Tangani error Tripay timeout/down
			if strings.Contains(err.Error(), "timeout") {
				return errs.ErrExternalTimeout
			}
			if strings.Contains(err.Error(), "connection error") {
				return errs.ErrExternalDown
			}
			return fmt.Errorf("create Tripay payment failed: %w", err)
		}

		booking.ReferenceNumber = &payment.Reference
		if err := uc.BookingRepository.Update(ctx, tx, booking); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to update booking with reference number: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	return payment, nil
}

func (uc *PaymentUsecase) HandleCallback(ctx context.Context, request *domain.Callback) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByOrderID(ctx, tx, request.MerchantRef)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}
		tickets, err := uc.TicketRepository.FindByBookingID(ctx, tx, booking.ID)
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
			if err := uc.HandleSuccessfulPayment(ctx, tx, booking, tickets); err != nil {
				return err
			}

		case "FAILED", "EXPIRED", "REFUND":
			// Payment unsuccessful
			if err := uc.HandleUnsuccessfulPayment(ctx, tx, booking, tickets, request.Status); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown payment status: %s", request.Status)
		}
		return nil
	})
}

func (uc *PaymentUsecase) HandleSuccessfulPayment(ctx context.Context, tx gotann.Connection, booking *domain.Booking, tickets []*domain.Ticket) error {
	// Send confirmation email
	subject := "Your Booking is Confirmed"
	htmlBody := templates.BookingSuccessEmail(booking, tickets)

	uc.Mailer.SendAsync(booking.Email, subject, htmlBody)
	return nil

}

func (uc *PaymentUsecase) HandleUnsuccessfulPayment(ctx context.Context, tx gotann.Connection, booking *domain.Booking, tickets []*domain.Ticket, status string) error {
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
		quota, err := uc.QuotaRepository.FindByScheduleIDAndClassID(ctx, tx, booking.ScheduleID, ticket.ClassID)
		if err == nil && quota != nil {
			// Count tickets for this class
			count := 0
			for _, t := range tickets {
				if t != nil && t.ClassID == ticket.ClassID {
					count++
				}
			}
			quota.Quota += count
			if err := uc.QuotaRepository.Update(ctx, tx, quota); err != nil {
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
		htmlBody = templates.BookingFailedEmail(booking, "Payment processing failed")
	case "EXPIRED":
		htmlBody = templates.BookingFailedEmail(booking, "Payment time expired")
	case "CANCELLED", "REFUND":
		htmlBody = templates.BookingFailedEmail(booking, "Payment was cancelled or refunded")
	}

	uc.Mailer.SendAsync(booking.Email, subject, htmlBody)
	return nil
}
