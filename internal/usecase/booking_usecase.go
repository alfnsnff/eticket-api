package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type BookingUsecase struct {
	Transactor        transact.Transactor
	BookingRepository domain.BookingRepository
	QuotaRepository   domain.QuotaRepository
}

func NewBookingUsecase(
	transactor transact.Transactor,
	booking_repository domain.BookingRepository,
	quota_repository domain.QuotaRepository,
) *BookingUsecase {
	return &BookingUsecase{
		Transactor:        transactor,
		BookingRepository: booking_repository,
		QuotaRepository:   quota_repository,
	}
}

func (uc *BookingUsecase) CreateBooking(ctx context.Context, e *domain.Booking) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking := &domain.Booking{
			OrderID:         e.OrderID,
			ScheduleID:      e.ScheduleID,
			IDType:          e.IDType,
			IDNumber:        e.IDNumber,
			CustomerName:    e.CustomerName,
			PhoneNumber:     e.PhoneNumber,
			Email:           e.Email,
			ReferenceNumber: e.ReferenceNumber,
		}
		if err := uc.BookingRepository.Insert(ctx, tx, booking); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create booking: %w", err)
		}
		return nil
	})
}

func (uc *BookingUsecase) ListBookings(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Booking, int, error) {
	var err error
	var total int64
	var bookings []*domain.Booking
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.BookingRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count bookings: %w", err)
		}

		bookings, err = uc.BookingRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all bookings: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, int(total), nil
}

func (uc *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var err error
	var booking *domain.Booking
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err = uc.BookingRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound

		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get booking by id: %w", err)
	}

	return booking, nil
}

func (uc *BookingUsecase) GetBookingByOrderID(ctx context.Context, orderID string) (*domain.Booking, error) {
	var err error
	var booking *domain.Booking
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err = uc.BookingRepository.FindByOrderID(ctx, tx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get booking by order ID: %w", err)
	}

	return booking, nil
}

func (uc *BookingUsecase) UpdateBooking(ctx context.Context, e *domain.Booking) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		booking.OrderID = e.OrderID
		booking.ScheduleID = e.ScheduleID
		booking.IDType = e.IDType
		booking.IDNumber = e.IDNumber
		booking.CustomerName = e.CustomerName
		booking.PhoneNumber = e.PhoneNumber
		booking.Email = e.Email
		booking.ReferenceNumber = e.ReferenceNumber

		if err := uc.BookingRepository.Update(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to update booking: %w", err)
		}
		return nil
	})
}

func (uc *BookingUsecase) RefundBooking(ctx context.Context, orderId string, email, IdNumber string) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByOrderID(ctx, tx, orderId)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		// Validate customer data without revealing which field is wrong
		if booking.Email != email ||
			booking.IDNumber != IdNumber {
			return fmt.Errorf("customer information does not match")
		}

		// Check if booking can be refunded
		if booking.Status != "PAID" {
			return fmt.Errorf("booking is not eligible for refund")
		}

		// Check if tickets exist
		if len(booking.Tickets) == 0 {
			return fmt.Errorf("no tickets found for this booking")
		}

		// Update booking status to refunded
		booking.Status = "REFUND"
		if err := uc.BookingRepository.Update(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to process refund")
		}

		restored := make(map[uint]bool)
		for _, ticket := range booking.Tickets {
			// Skip if already processed this class
			if restored[ticket.ClassID] {
				continue
			}

			quota, err := uc.QuotaRepository.FindByScheduleIDAndClassID(ctx, tx, booking.ScheduleID, ticket.ClassID)
			if err == nil && quota != nil {
				// Count tickets for this class
				count := 0
				for _, t := range booking.Tickets {
					if t.ClassID == ticket.ClassID {
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

		return nil
	})
}

func (uc *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		booking, err := uc.BookingRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get booking: %w", err)
		}
		if booking == nil {
			return errs.ErrNotFound
		}

		if err := uc.BookingRepository.Delete(ctx, tx, booking); err != nil {
			return fmt.Errorf("failed to delete booking: %w", err)
		}
		return nil
	})
}
