package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingUsecase struct {
	DB                *gorm.DB
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewBookingUsecase(db *gorm.DB,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
	session_repository *repository.SessionRepository,

) *BookingUsecase {
	return &BookingUsecase{
		DB:                db,
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

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		return b.BookingRepository.Create(tx, booking)
	})
}

func (b *BookingUsecase) GetAllBookings(ctx context.Context) ([]*model.ReadBookingResponse, error) {
	bookings := []*entity.Booking{}

	err := tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		var err error
		bookings, err = b.BookingRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.BookingMapper.ToModels(bookings), nil
}

func (b *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	booking := new(entity.Booking)

	err := tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
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

func (b *BookingUsecase) UpdateBooking(ctx context.Context, id uint, request *model.UpdateBookingRequest) error {
	booking := mapper.BookingMapper.FromUpdate(request)
	booking.ID = id

	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CustomerName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
		return b.BookingRepository.Update(tx, booking)
	})
}

// Deletebooking deletes a booking by its ID
func (b *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {

	return tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {
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

func (b *BookingUsecase) ConfirmBooking(ctx context.Context, request *model.ConfirmBookingRequest) (*model.ConfirmBookingResponse, error) {
	if err := HelperValidateConfirmRequest(request); err != nil {
		return nil, err
	}

	var confirmedBooking *entity.Booking
	var confirmedTicketIDs []uint

	err := tx.Execute(ctx, b.DB, func(tx *gorm.DB) error {

		session, err := b.SessionRepository.GetByUUIDWithLock(tx, request.SessionID, true)
		if err != nil {
			return fmt.Errorf("failed to retrieve claim session %s within transaction: %w", request.SessionID, err)
		}
		if session == nil {
			return errors.New("claim session not found")
		}

		tickets, err := b.TicketRepository.FindManyBySessionID(tx, session.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve tickets: %w", err)
		}

		now := time.Now()
		ticketsToUpdate, total, scheduleID, err := b.HelperValidateAndPrepareTickets(tickets, now)
		if err != nil {
			return err
		}

		booking := HelperBuildBooking(request, scheduleID, now, total)
		if err := b.BookingRepository.Create(tx, booking); err != nil {
			return fmt.Errorf("failed to create booking: %w", err)
		}
		confirmedBooking = booking

		confirmedTicketIDs = HelperUpdateTicketsWithBooking(ticketsToUpdate, booking.ID)

		if err := b.TicketRepository.UpdateBulk(tx, ticketsToUpdate); err != nil {
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

func HelperValidateConfirmRequest(request *model.ConfirmBookingRequest) error {
	if request.Name == "" || request.IDType == "" || request.IDNumber == "" ||
		request.PhoneNumber == "" || request.Email == "" || request.BirthDate.IsZero() {
		return errors.New("invalid request: missing required fields")
	}
	return nil
}

func (b *BookingUsecase) HelperValidateAndPrepareTickets(
	tickets []*entity.Ticket,
	now time.Time,
) ([]*entity.Ticket, float32, uint, error) {
	var total float32
	var scheduleID uint
	for i, ticket := range tickets {
		if ticket.Status != "pending_payment" {
			return nil, 0, 0, fmt.Errorf("ticket %d not in pending_payment state", ticket.ID)
		}
		// if ticket.ExpiresAt.Before(now) {
		// 	return nil, 0, 0, fmt.Errorf("ticket %d expired before confirmation", ticket.ID)
		// }
		if i == 0 {
			scheduleID = ticket.ScheduleID
		}
		ticket.Status = "confirmed"
		ticket.ClaimSessionID = nil
		ticket.BookedAt = &now
		total += ticket.Price
	}
	return tickets, total, scheduleID, nil
}

func HelperBuildBooking(request *model.ConfirmBookingRequest, scheduleID uint, now time.Time, total float32) *entity.Booking {
	return &entity.Booking{
		ScheduleID:   scheduleID,
		CustomerName: request.Name,
		IDType:       request.IDType,
		IDNumber:     request.IDNumber,
		PhoneNumber:  request.PhoneNumber,
		Email:        request.Email,
		BirthDate:    request.BirthDate,
		BookedAt:     now,
		TotalPrice:   total,
		Status:       "completed",
		// PaymentIntentID: request.PaymentIntentID, // Uncomment if needed
	}
}

func HelperUpdateTicketsWithBooking(tickets []*entity.Ticket, bookingID uint) []uint {
	ids := make([]uint, len(tickets))
	for i, t := range tickets {
		t.BookingID = &bookingID
		ids[i] = t.ID
	}
	return ids
}
