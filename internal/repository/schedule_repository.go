package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{DB: db}
}

func (r *ScheduleRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.Schedule{}).Count(&total)
	return total, result.Error
}

func (r *ScheduleRepository) Insert(ctx context.Context, conn gotann.Connection, schedule *domain.Schedule) error {
	result := conn.Create(schedule)
	return result.Error
}

// Add bulk operations for consistency
func (r *ScheduleRepository) InsertBulk(ctx context.Context, conn gotann.Connection, schedules []*domain.Schedule) error {
	result := conn.Create(&schedules)
	return result.Error
}

func (r *ScheduleRepository) Update(ctx context.Context, conn gotann.Connection, schedule *domain.Schedule) error {
	result := conn.Save(schedule)
	return result.Error
}

func (r *ScheduleRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, schedules []*domain.Schedule) error {
	result := conn.Save(&schedules)
	return result.Error
}

func (r *ScheduleRepository) Delete(ctx context.Context, conn gotann.Connection, schedule *domain.Schedule) error {
	result := conn.Select(clause.Associations).Delete(schedule)
	return result.Error
}

func (r *ScheduleRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	query := conn.Model(&domain.Schedule{}).Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Preload("Ship").
		Preload("Quotas").
		Preload("Quotas.Class")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("schedule_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.Schedule, error) {
	schedule := new(domain.Schedule)
	result := conn.
		Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Preload("Ship").
		Preload("Quotas").
		Preload("Quotas.Class").
		First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return schedule, result.Error
}

func (r *ScheduleRepository) FindByStatus(ctx context.Context, conn gotann.Connection, status string) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	result := conn.
		Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Preload("Ship").
		Preload("Quotas").
		Preload("Quotas.Class").
		Where("status = ?", status).
		Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

func (r *ScheduleRepository) FindActiveSchedules(ctx context.Context, conn gotann.Connection) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	result := conn.
		Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Preload("Ship").
		Preload("Quotas").
		Preload("Quotas.Class").
		Where("departure_datetime > ?", time.Now()).
		Limit(7).
		Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}
