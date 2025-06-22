package repository

import (
	"errors"
	"eticket-api/internal/entity"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ScheduleRepository struct{}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

func (ar *ScheduleRepository) Create(db *gorm.DB, schedule *entity.Schedule) error {
	result := db.Create(schedule)
	return result.Error
}

func (ar *ScheduleRepository) Update(db *gorm.DB, schedule *entity.Schedule) error {
	result := db.Save(schedule)
	return result.Error
}

func (ar *ScheduleRepository) Delete(db *gorm.DB, schedule *entity.Schedule) error {
	result := db.Select(clause.Associations).Delete(schedule)
	return result.Error
}

func (scr *ScheduleRepository) Count(db *gorm.DB) (int64, error) {
	schedules := []*entity.Schedule{}
	var total int64
	result := db.Find(&schedules).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (scr *ScheduleRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.Schedule, error) {
	schedules := []*entity.Schedule{}

	query := db.Preload("Route").
		Preload("Route.DepartureHarbor").
		Preload("Route.ArrivalHarbor").
		Preload("Ship")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("schedule_id ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&schedules).Error
	return schedules, err
}

func (scr *ScheduleRepository) GetByID(db *gorm.DB, id uint) (*entity.Schedule, error) {
	schedule := new(entity.Schedule)
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return schedule, result.Error
}

func (scr *ScheduleRepository) GetAllScheduled(db *gorm.DB) ([]*entity.Schedule, error) {
	schedules := []*entity.Schedule{}

	// Corrected line: Use "?" as a placeholder and pass the string value as a parameter
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Where("status = ?", "scheduled").Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}

func (scr *ScheduleRepository) GetActiveSchedule(db *gorm.DB) ([]*entity.Schedule, error) {
	schedules := []*entity.Schedule{}

	// Corrected line: Use "?" as a placeholder and pass the string value as a parameter
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Where("departure_datetime > ?", time.Now()).Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}
