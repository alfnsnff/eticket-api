package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository struct{}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{}
}

func (br *BookingRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Booking{}).Count(&total)
	return total, result.Error
}

func (ar *BookingRepository) Insert(db *gorm.DB, booking *domain.Booking) error {
	result := db.Create(booking)
	return result.Error
}

func (br *BookingRepository) InsertBulk(db *gorm.DB, bookings []*domain.Booking) error {
	result := db.Create(&bookings)
	return result.Error
}

func (br *BookingRepository) Update(db *gorm.DB, booking *domain.Booking) error {
	result := db.Save(booking)
	return result.Error
}

func (br *BookingRepository) UpdateBulk(db *gorm.DB, bookings []*domain.Booking) error {
	result := db.Save(&bookings)
	return result.Error
}

func (br *BookingRepository) Delete(db *gorm.DB, booking *domain.Booking) error {
	result := db.Select(clause.Associations).Delete(booking)
	return result.Error
}

func (br *BookingRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Booking, error) {
	bookings := []*domain.Booking{}
	query := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("order_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}

func (br *BookingRepository) FindByID(db *gorm.DB, id uint) (*domain.Booking, error) {
	booking := new(domain.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		First(&booking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}

func (br *BookingRepository) FindByOrderID(db *gorm.DB, id string) (*domain.Booking, error) {
	booking := new(domain.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Where("order_id = ?", id).First(&booking)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}
