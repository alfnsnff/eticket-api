package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

func (r *BookingRepository) Create(db *gorm.DB, booking *entities.Booking) error {
	result := db.Create(booking) // GORM automatically assigns ID after insert
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetAll retrieves all bookings from the database
func (r *BookingRepository) GetAll(db *gorm.DB) ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	result := db.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Schedule").
		Preload("Tickets").
		Preload("Tickets.Price.ShipClass.Class").
		Preload("Tickets.Price.ShipClass").
		Preload("Tickets.Price").
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
		Preload("Tickets.Price.ShipClass.Class").
		Preload("Tickets.Price.ShipClass").
		Preload("Tickets.Price").
		First(&booking, id) // Fetches the booking by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no booking is found
	}
	return &booking, result.Error
}

// Update modifies an existing booking in the database
func (r *BookingRepository) Update(db *gorm.DB, booking *entities.Booking) error {
	// Uses Gorm's Save method to update the booking
	result := db.Save(booking)
	return result.Error
}

// Delete removes a booking from the database by its ID
func (r *BookingRepository) Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&entities.Booking{}, id) // Deletes the booking by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no booking found to delete") // Custom error for non-existent ID
	}
	return nil
}
