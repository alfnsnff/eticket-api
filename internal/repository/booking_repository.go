package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) entities.BookingRepositoryInterface {
	return &BookingRepository{DB: db}
}

// Create inserts a new booking into the database
func (r *BookingRepository) Create(booking *entities.Booking) error {
	result := r.DB.Create(booking)
	return result.Error
}

// GetAll retrieves all bookings from the database
func (r *BookingRepository) GetAll() ([]*entities.Booking, error) {
	var bookings []*entities.Booking
	result := r.DB.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Schedule").Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

// GetByID retrieves a booking by its ID
func (r *BookingRepository) GetByID(id uint) (*entities.Booking, error) {
	var booking entities.Booking
	result := r.DB.Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Schedule").First(&booking, id) // Fetches the booking by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no booking is found
	}
	return &booking, result.Error
}

// Update modifies an existing booking in the database
func (r *BookingRepository) Update(booking *entities.Booking) error {
	// Uses Gorm's Save method to update the booking
	result := r.DB.Save(booking)
	return result.Error
}

// Delete removes a booking from the database by its ID
func (r *BookingRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Booking{}, id) // Deletes the booking by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no booking found to delete") // Custom error for non-existent ID
	}
	return nil
}
