package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{DB: db}
}

// Create inserts a new schedule into the database
func (r *ScheduleRepository) Create(schedule *entities.Schedule) error {
	result := r.DB.Create(schedule)
	return result.Error
}

// GetAll retrieves all schedules from the database
func (r *ScheduleRepository) GetAll() ([]*entities.Schedule, error) {
	var schedules []*entities.Schedule
	result := r.DB.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

// GetByID implements entities.ScheduleRepositoryInterface.
func (r *ScheduleRepository) GetByID(id uint) (*entities.Schedule, error) {
	var schedule entities.Schedule
	result := r.DB.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").First(&schedule, id) // Fetches the schedule by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no schedule is found
	}
	return &schedule, result.Error
}

func (r *ScheduleRepository) Search(routeID uint, date time.Time, shipID *uint) (*entities.Schedule, error) {
	var schedule *entities.Schedule

	query := r.DB.
		Preload("Route.DepartureHarbor").
		Preload("Route.ArrivalHarbor").
		Preload("Ship").
		Where("route_id = ?", routeID).
		Where("DATE(datetime) = ?", date.Format("2006-01-02"))

	if shipID != nil {
		query = query.Where("ship_id = ?", *shipID)
	}

	result := query.Find(&schedule)
	return schedule, result.Error
}

// Update modifies an existing schedule in the database
func (r *ScheduleRepository) Update(schedule *entities.Schedule) error {
	// Uses Gorm's Save method to update the schedule
	result := r.DB.Save(schedule)
	return result.Error
}

// Delete removes a schedule from the database by its ID
func (r *ScheduleRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Schedule{}, id) // Deletes the schedule by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no schedule found to delete") // Custom error for non-existent ID
	}
	return nil
}
