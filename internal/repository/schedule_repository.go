package repository

import (
	"errors"
	"eticket-api/internal/domain"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) domain.ScheduleRepositoryInterface {
	return &ScheduleRepository{DB: db}
}

// Create inserts a new schedule into the database
func (r *ScheduleRepository) Create(schedule *domain.Schedule) error {
	result := r.DB.Create(schedule)
	return result.Error
}

// GetAll retrieves all schedules from the database
func (r *ScheduleRepository) GetAll() ([]*domain.Schedule, error) {
	var schedules []*domain.Schedule
	result := r.DB.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

// GetByID retrieves a schedule by its ID
func (r *ScheduleRepository) GetByID(id uint) (*domain.Schedule, error) {
	var schedule domain.Schedule
	result := r.DB.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").First(&schedule, id) // Fetches the schedule by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no schedule is found
	}
	return &schedule, result.Error
}

// Update modifies an existing schedule in the database
func (r *ScheduleRepository) Update(schedule *domain.Schedule) error {
	// Uses Gorm's Save method to update the schedule
	result := r.DB.Save(schedule)
	return result.Error
}

// Delete removes a schedule from the database by its ID
func (r *ScheduleRepository) Delete(id uint) error {
	result := r.DB.Delete(&domain.Schedule{}, id) // Deletes the schedule by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no schedule found to delete") // Custom error for non-existent ID
	}
	return nil
}
