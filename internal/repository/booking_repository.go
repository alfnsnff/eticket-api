package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"

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

func (br *BookingRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Booking, error) {
	bookings := []*entity.Booking{}

	query := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("order_id ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}

func (br *BookingRepository) GetByID(db *gorm.DB, id uint) (*entity.Booking, error) {
	booking := new(entity.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		First(&booking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}

func (br *BookingRepository) GetByOrderID(db *gorm.DB, id string) (*entity.Booking, error) {
	booking := new(entity.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Where("order_id = ?", id).Find(&booking)
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
	return result.Error
}

func (r *BookingRepository) UpdateReferenceNumber(tx *gorm.DB, bookingID uint, reference *string) error {
	return tx.Model(&entity.Booking{}).
		Where("id = ?", bookingID).
		Update("reference_number", reference).Error
}
