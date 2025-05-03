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
