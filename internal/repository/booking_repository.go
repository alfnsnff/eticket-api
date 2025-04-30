package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type BookingRepository struct {
	Repository[entities.Booking]
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

// GetAll retrieves all bookings from the database
func (r *BookingRepository) GetAll(db *gorm.DB) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	result := db.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Schedule").
		Preload("Tickets").
		Preload("Tickets.Fare.Manifest.Class").
		Preload("Tickets.Fare.Manifest").
		Preload("Tickets.Fare").
		Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

// GetByID retrieves a booking by its ID
func (r *BookingRepository) GetByID(db *gorm.DB, id uint) (*entities.Booking, error) {
	var booking entities.Booking
	result := db.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Schedule").
		Preload("Tickets").
		Preload("Tickets.Fare.Manifest.Class").
		Preload("Tickets.Fare.Manifest").
		Preload("Tickets.Fare").
		First(&booking, id) // Fetches the booking by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no booking is found
	}
	return &booking, result.Error
}
