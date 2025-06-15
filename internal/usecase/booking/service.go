package booking

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingUsecase struct {
	Tx                *tx.TxManager
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewBookingUsecase(
	tx *tx.TxManager,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
	session_repository *repository.SessionRepository,

) *BookingUsecase {
	return &BookingUsecase{
		Tx:                tx,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
		SessionRepository: session_repository,
	}
}

func (b *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	booking := mapper.BookingMapper.FromWrite(request)

	if booking.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return b.BookingRepository.Create(tx, booking)
	})
}

func (b *BookingUsecase) GetAllBookings(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadBookingResponse, int, error) {
	bookings := []*entity.Booking{}
	var total int64
	err := b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = b.BookingRepository.Count(tx)
		if err != nil {
			return err
		}
		bookings, err = b.BookingRepository.GetAll(tx, limit, offset, sort, search)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.BookingMapper.ToModels(bookings), int(total), nil
}

func (b *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	booking := new(entity.Booking)

	err := b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		booking, err = b.BookingRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.BookingMapper.ToModel(booking), nil
}

func (b *BookingUsecase) GetBookingByOrderID(ctx context.Context, id string) (*model.ReadBookingResponse, error) {
	booking := new(entity.Booking)

	err := b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		booking, err = b.BookingRepository.GetByOrderID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.BookingMapper.ToModel(booking), nil
}

func (b *BookingUsecase) UpdateBooking(ctx context.Context, id uint, request *model.UpdateBookingRequest) error {
	booking := mapper.BookingMapper.FromUpdate(request)
	booking.ID = id

	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return b.BookingRepository.Update(tx, booking)
	})
}

// Deletebooking deletes a booking by its ID
func (b *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {

	return b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		booking, err := b.BookingRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return errors.New("route not found")
		}
		return b.BookingRepository.Delete(tx, booking)
	})

}

func (b *BookingUsecase) PaidConfirm(ctx context.Context, id uint) error {
	return b.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return b.BookingRepository.PaidConfirm(tx, id)
	})
}

func (b *BookingUsecase) ConfirmBooking(ctx context.Context, sessionID string) (*model.ConfirmBookingResponse, error) {
	// Inlined HelperValidateConfirmRequest logic
	var confirmedBooking *entity.Booking
	var confirmedTicketIDs []uint

	err := b.Tx.Execute(ctx, func(tx *gorm.DB) error {

		session, err := b.SessionRepository.GetByUUIDWithLock(tx, sessionID, true)
		if err != nil {
			return fmt.Errorf("failed to retrieve claim session %s within transaction: %w", sessionID, err)
		}
		if session == nil {
			return errors.New("claim session not found")
		}

		tickets, err := b.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}

		var total float32
		now := time.Now()
		for _, ticket := range tickets {
			if ticket.Status != "pending_payment" {
				return fmt.Errorf("ticket %d not in pending_payment state", ticket.ID)
			}
			ticket.Status = "paid_confirmed"
			ticket.ClaimSessionID = nil
			ticket.BookedAt = &now
			total += ticket.Price
		}

		// Inlined HelperBuildBooking logic
		booking := &entity.Booking{
			Status: "paid_comfirmed",
		}

		if err := b.BookingRepository.Update(tx, booking); err != nil {
			return fmt.Errorf("failed to update booking: %w", err)
		}
		confirmedBooking = booking

		// Inlined HelperUpdateTicketsWithBooking logic
		ids := make([]uint, len(tickets))
		for i, t := range tickets {
			t.BookingID = &booking.ID
			ids[i] = t.ID
		}
		confirmedTicketIDs = ids

		if err := b.TicketRepository.UpdateBulk(tx, tickets); err != nil {
			return fmt.Errorf("failed to update tickets: %w", err)
		}
		return nil

	})

	if err != nil {
		return nil, fmt.Errorf("confirm book transaction failed: %w", err)
	}

	return &model.ConfirmBookingResponse{
		BookingID:          confirmedBooking.ID,
		BookingStatus:      confirmedBooking.Status,
		ConfirmedTicketIDs: confirmedTicketIDs,
	}, nil
}
