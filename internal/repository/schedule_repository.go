package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepository struct{}

// FindAllScheduled implements domain.ScheduleRepository.
func (sr *ScheduleRepository) FindAllScheduled(db *gorm.DB) ([]*domain.Schedule, error) {
	panic("unimplemented")
}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

func (sr *ScheduleRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.Schedule{}).Count(&total)
	return total, result.Error
}

func (sr *ScheduleRepository) Insert(db *gorm.DB, schedule *domain.Schedule) error {
	result := db.Create(schedule)
	return result.Error
}

func (sr *ScheduleRepository) Inserts(conn *gotann.Connection, schedule *domain.Schedule) error {
	db, err := gotann.GetGormDB(conn)
	if err != nil {
		return fmt.Errorf("failed to get GORM DB: %w", err)
	}
	result := db.Create(schedule)
	return result.Error
}

// Add bulk operations for consistency
func (sr *ScheduleRepository) InsertBulk(db *gorm.DB, schedules []*domain.Schedule) error {
	result := db.Create(&schedules)
	return result.Error
}

func (sr *ScheduleRepository) Update(db *gorm.DB, schedule *domain.Schedule) error {
	result := db.Save(schedule)
	return result.Error
}

func (sr *ScheduleRepository) UpdateBulk(db *gorm.DB, schedules []*domain.Schedule) error {
	result := db.Save(&schedules)
	return result.Error
}

func (sr *ScheduleRepository) Delete(db *gorm.DB, schedule *domain.Schedule) error {
	result := db.Select(clause.Associations).Delete(schedule)
	return result.Error
}

func (sr *ScheduleRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	query := db.Preload("DepartureHarbor").
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

func (sr *ScheduleRepository) FindByID(db *gorm.DB, id uint) (*domain.Schedule, error) {
	schedule := new(domain.Schedule)
	result := db.
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

func (sr *ScheduleRepository) FindByStatus(db *gorm.DB, status string) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	result := db.
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

func (sr *ScheduleRepository) FindActiveSchedule(db *gorm.DB) ([]*domain.Schedule, error) {
	schedules := []*domain.Schedule{}
	result := db.
		Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Preload("Ship").
		Preload("Quotas").
		Preload("Quotas.Class").
		Where("departure_datetime > ?", time.Now()).
		Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}
