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

func (br *BookingRepository) Count(db *gorm.DB) (int64, error) {
	bookings := []*entity.Booking{}
	var total int64
	result := db.Find(&bookings).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (br *BookingRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Booking, error) {
	bookings := []*entity.Booking{}
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Limit(limit).Offset(offset).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

func (br *BookingRepository) GetByID(db *gorm.DB, id uint) (*entity.Booking, error) {
	booking := new(entity.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		First(&booking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}

func (br *BookingRepository) PaidConfirm(db *gorm.DB, id uint) error {
	booking := new(entity.Booking)
	result := db.First(&booking, id).Update("status", "paid")
	if result.Error != nil {
		return result.Error
	}
	return nil
}
