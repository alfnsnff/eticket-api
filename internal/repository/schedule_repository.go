package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	Repository[entities.Schedule]
}

func NewScheduleRepository() *ScheduleRepository {
	return &ScheduleRepository{}
}

// GetAll retrieves all schedules from the database
func (r *ScheduleRepository) GetAll(db *gorm.DB) ([]*entities.Schedule, error) {
	var schedules []*entities.Schedule
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

// GetByID implements entities.ScheduleRepositoryInterface.
func (r *ScheduleRepository) GetByID(db *gorm.DB, id uint) (*entities.Schedule, error) {
	var schedule entities.Schedule
	result := db.Preload("Route").Preload("Route.DepartureHarbor").Preload("Route.ArrivalHarbor").Preload("Ship").First(&schedule, id) // Fetches the schedule by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no schedule is found
	}
	return &schedule, result.Error
}

func (r *ScheduleRepository) Search(db *gorm.DB, routeID uint, date time.Time, shipID *uint) (*entities.Schedule, error) {
	var schedule *entities.Schedule

	query := db.
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
