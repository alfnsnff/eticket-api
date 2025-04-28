package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
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
	bookingRepository *repository.BookingRepository,
	ticketRepository *repository.TicketRepository,
) *BookingUsecase {
	return &BookingUsecase{
		DB:                db,
		BookingRepository: bookingRepository,
		TicketRepository:  ticketRepository,
	}
}

// Createbooking validates and creates a new booking
func (s *BookingUsecase) CreateBooking(ctx context.Context, booking *entities.Booking) error {
	// booking, _ := dto.ToBookingEntity(bookingCreate)

	if booking.CusName == "" {
		return fmt.Errorf("customer name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.BookingRepository.Create(txDB, booking)
	})

	// return s.BookingRepository.Create(booking)
}

// func (s *BookingUsecase) CreateBookingWithTickets(ctx context.Context, booking *entities.Booking, tickets *[]entities.Ticket) error {
// 	// booking, tickets := dto.ToBookingEntity(bookingCreate)

// 	// Validate booking
// 	if booking.CusName == "" {
// 		return fmt.Errorf("customer name cannot be empty")
// 	}

// 	// // Create the booking first
// 	// err := s.BookingRepository.Create(booking)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // Loop through tickets and create each one
// 	// for i := range *tickets {
// 	// 	(*tickets)[i].BookingID = booking.ID
// 	// 	// (*tickets)[i].ScheduleID = booking.Schedule.ID

// 	// 	err := s.TicketRepository.Create(&(*tickets)[i])
// 	// 	if err != nil {
// 	// 		return fmt.Errorf("failed to create ticket: %v", err)
// 	// 	}
// 	// }

// 	return s.BookingRepository.Create(booking)
// }

func (s *BookingUsecase) GetAllBookings(ctx context.Context) ([]*entities.Booking, error) {
	var bookings []*entities.Booking

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		bookings, err = s.BookingRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return bookings, nil
}

func (s *BookingUsecase) GetBookingByID(ctx context.Context, id uint) (*entities.Booking, error) {
	var booking *entities.Booking

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		booking, err = s.BookingRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	return booking, nil
}

// Updatebooking updates an existing booking
func (s *BookingUsecase) UpdateBooking(ctx context.Context, id uint, booking *entities.Booking) error {
	booking.ID = id

	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CusName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.BookingRepository.Update(txDB, booking)
	})
}

// Deletebooking deletes a booking by its ID
func (s *BookingUsecase) DeleteBooking(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		booking, err := s.BookingRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if booking == nil {
			return errors.New("route not found")
		}
		return s.BookingRepository.Delete(txDB, id)
	})
}
