package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	Repository[entity.Schedule]
}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

func (scr *ScheduleRepository) GetAll(db *gorm.DB) ([]*entity.Schedule, error) {
	schedules := []*entity.Schedule{}
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
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
