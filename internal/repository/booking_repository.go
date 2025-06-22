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

func (ar *BookingRepository) Create(db *gorm.DB, booking *domain.Booking) error {
	result := db.Create(booking)
	return result.Error
}

func (ar *BookingRepository) Update(db *gorm.DB, booking *domain.Booking) error {
	result := db.Save(booking)
	return result.Error
}

func (ar *BookingRepository) Delete(db *gorm.DB, booking *domain.Booking) error {
	result := db.Select(clause.Associations).Delete(booking)
	return result.Error
}

func (br *BookingRepository) Count(db *gorm.DB) (int64, error) {
	bookings := []*domain.Booking{}
	var total int64
	result := db.Find(&bookings).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (br *BookingRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Booking, error) {
	bookings := []*domain.Booking{}

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

func (br *BookingRepository) GetByID(db *gorm.DB, id uint) (*domain.Booking, error) {
	booking := new(domain.Booking)
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

// ✅ FIXED - Use First() for single record
func (br *BookingRepository) GetByOrderID(db *gorm.DB, id string) (*domain.Booking, error) {
	booking := new(domain.Booking)
	result := db.Preload("Tickets").
		Preload("Tickets.Class").
		Preload("Schedule").
		Preload("Schedule.Ship").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Where("order_id = ?", id).First(&booking) // ✅ Use First()

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return booking, result.Error
}

func (br *BookingRepository) PaidConfirm(db *gorm.DB, id uint) error {
	booking := new(domain.Booking)
	result := db.First(&booking, id).Update("status", "paid")
	if result.Error != nil {
		return result.Error
	}
	return result.Error
}

func (r *BookingRepository) UpdateReferenceNumber(tx *gorm.DB, bookingID uint, reference *string) error {
	return tx.Model(&domain.Booking{}).
		Where("id = ?", bookingID).
		Update("reference_number", reference).Error
}
