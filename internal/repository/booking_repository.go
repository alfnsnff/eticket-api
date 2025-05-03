package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type BookingRepository struct {
	Repository[entity.Booking]
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

func (br *BookingRepository) GetAll(db *gorm.DB) ([]*entity.Booking, error) {
	bookings := []*entity.Booking{}
	result := db.Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

func (br *BookingRepository) GetByID(db *gorm.DB, id uint) (*entity.Booking, error) {
	booking := new(entity.Booking)
	result := db.First(&booking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}
