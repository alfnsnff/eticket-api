package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"fmt"
)

type BookingUsecase struct {
	BookingRepository entities.BookingRepositoryInterface
}

func NewBookingUsecase(bookingRepository entities.BookingRepositoryInterface) BookingUsecase {
	return BookingUsecase{BookingRepository: bookingRepository}
}

// Createbooking validates and creates a new booking
func (s *BookingUsecase) CreateBooking(booking *entities.Booking) error {
	if booking.CusName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}
	return s.BookingRepository.Create(booking)
}

// GetAllbookinges retrieves all bookings
func (s *BookingUsecase) GetAllBookings() ([]*entities.Booking, error) {
	return s.BookingRepository.GetAll()
}

// GetbookingByID retrieves a booking by its ID
func (s *BookingUsecase) GetBookingByID(id uint) (*entities.Booking, error) {
	booking, err := s.BookingRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	return booking, nil
}

// Updatebooking updates an existing booking
func (s *BookingUsecase) UpdateBooking(booking *entities.Booking) error {
	if booking.ID == 0 {
		return fmt.Errorf("booking ID cannot be zero")
	}
	if booking.CusName == "" {
		return fmt.Errorf("booking name cannot be empty")
	}
	return s.BookingRepository.Update(booking)
}

// Deletebooking deletes a booking by its ID
func (s *BookingUsecase) DeleteBooking(id uint) error {
	booking, err := s.BookingRepository.GetByID(id)
	if err != nil {
		return err
	}
	if booking == nil {
		return errors.New("booking not found")
	}
	return s.BookingRepository.Delete(id)
}
