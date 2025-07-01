package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository struct {
	DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{DB: db}
}

func (r *BookingRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Booking{}).Count(&total)
	return total, result.Error
}

func (r *BookingRepository) Insert(ctx context.Context, conn gotann.Connection, booking *domain.Booking) error {
	result := conn.Create(booking)
	return result.Error
}

func (r *BookingRepository) InsertBulk(ctx context.Context, conn gotann.Connection, bookings []*domain.Booking) error {
	result := conn.Create(&bookings)
	return result.Error
}

func (r *BookingRepository) Update(ctx context.Context, conn gotann.Connection, booking *domain.Booking) error {
	result := conn.Save(booking)
	return result.Error
}

func (r *BookingRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, bookings []*domain.Booking) error {
	result := conn.Save(&bookings)
	return result.Error
}

func (r *BookingRepository) Delete(ctx context.Context, conn gotann.Connection, booking *domain.Booking) error {
	result := conn.Select(clause.Associations).Delete(booking)
	return result.Error
}

func (r *BookingRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Booking, error) {
	bookings := []*domain.Booking{}
	query := conn.Model(&domain.Booking{}).Preload("Tickets").
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

func (r *BookingRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Booking, error) {
	booking := new(domain.Booking)
	result := conn.Preload("Tickets").
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

func (r *BookingRepository) FindByOrderID(ctx context.Context, conn gotann.Connection, id string) (*domain.Booking, error) {
	booking := new(domain.Booking)
	result := conn.Preload("Tickets").
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
