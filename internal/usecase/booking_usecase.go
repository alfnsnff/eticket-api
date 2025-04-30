package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type BookingUsecase struct {
	DB                *gorm.DB
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
}

func NewBookingUsecase(db *gorm.DB,
	booking_repository *repository.BookingRepository,
	ticket_repository *repository.TicketRepository,
) *BookingUsecase {
	return &BookingUsecase{
		DB:                db,
		BookingRepository: booking_repository,
		TicketRepository:  ticket_repository,
	}
}

// Createbooking validates and creates a new booking
func (s *BookingUsecase) CreateBooking(ctx context.Context, request *model.WriteBookingRequest) error {
	booking := mapper.ToBookingEntity(request)

	if booking.CusName == "" {
		return fmt.Errorf("customer name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.BookingRepository.Create(tx, booking)
	})

	// return s.BookingRepository.Create(booking)
}

func (s *BookingUsecase) GetAllBookings(ctx context.Context) ([]*model.ReadBookingResponse, error) {
	bookings := []*entities.Booking{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		bookings, err = s.BookingRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ToBookingsModel(bookings), nil
}

func (s *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*model.ReadBookingResponse, error) {
	booking := new(entities.Booking)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		booking, err = s.BookingRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return mapper.ToBookingModel(booking), nil
}

// Updatebooking updates an existing booking
func (s *BookingUsecase) UpdateBooking(ctx context.Context, id uint, request *model.WriteBookingRequest) error {

	booking := mapper.ToBookingEntity(request)
	booking.ID = id

	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CusName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.BookingRepository.Update(tx, booking)
	})
}

// Deletebooking deletes a booking by its ID
func (s *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		booking, err := s.BookingRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return errors.New("route not found")
		}
		return s.BookingRepository.Delete(tx, booking)
	})
}
