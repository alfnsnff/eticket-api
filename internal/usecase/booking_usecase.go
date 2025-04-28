package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	"fmt"
)

type BookingUsecase struct {
	BookingRepository *repository.BookingRepository
	TicketRepository  *repository.TicketRepository
}

func NewBookingUsecase(
	bookingRepository *repository.BookingRepository,
	ticketRepository *repository.TicketRepository,
) BookingUsecase {
	return BookingUsecase{
		BookingRepository: bookingRepository,
		TicketRepository:  ticketRepository,
	}
}

// Createbooking validates and creates a new booking
func (s *BookingUsecase) CreateBooking(booking *entities.Booking) error {
	// booking, _ := dto.ToBookingEntity(bookingCreate)

	if booking.CusName == "" {
		return fmt.Errorf("customer name cannot be empty")
	}

	return s.BookingRepository.Create(booking)
}

func (s *BookingUsecase) CreateBookingWithTickets(booking *entities.Booking, tickets *[]entities.Ticket) error {
	// booking, tickets := dto.ToBookingEntity(bookingCreate)

	// Validate booking
	if booking.CusName == "" {
		return fmt.Errorf("customer name cannot be empty")
	}

	// // Create the booking first
	// err := s.BookingRepository.Create(booking)
	// if err != nil {
	// 	return err
	// }

	// // Loop through tickets and create each one
	// for i := range *tickets {
	// 	(*tickets)[i].BookingID = booking.ID
	// 	// (*tickets)[i].ScheduleID = booking.Schedule.ID

	// 	err := s.TicketRepository.Create(&(*tickets)[i])
	// 	if err != nil {
	// 		return fmt.Errorf("failed to create ticket: %v", err)
	// 	}
	// }

	return s.BookingRepository.Create(booking)
}

func (s *BookingUsecase) GetAllBookings() ([]*entities.Booking, error) {
	return s.BookingRepository.GetAll()
}

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
func (s *BookingUsecase) UpdateBooking(id uint, booking *entities.Booking) error {
	booking.ID = id

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
