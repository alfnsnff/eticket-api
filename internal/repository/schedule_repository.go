package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

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
	result := db.Find(&schedules)
	if result.Error != nil {
		return nil, result.Error
	}
	return schedules, nil
}

func (scr *ScheduleRepository) GetByID(db *gorm.DB, id uint) (*entity.Schedule, error) {
	schedule := new(entity.Schedule)
	result := db.First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return schedule, result.Error
}

// func (r *ScheduleRepository) SearchByRoute(db *gorm.DB, routeID uintAndTime) (*entity.Schedule, error) {
// 	var schedule *entity.Schedule

// 	query := db.
// 		Preload("Route.DepartureHarbor").
// 		Preload("Route.ArrivalHarbor").
// 		Preload("Ship").
// 		Where("route_id = ?", routeID).
// 		Where("DATE(datetime) = ?", date.Format("2006-01-02"))

// 	if shipID != nil {
// 		query = query.Where("ship_id = ?", *shipID)
// 	}

// 	result := query.Find(&schedule)
// 	return schedule, result.Error
// }

func (scr *ScheduleRepository) GetAllScheduled(db *gorm.DB) ([]*entity.Schedule, error) {
	schedules := []*entity.Schedule{}

	// Corrected line: Use "?" as a placeholder and pass the string value as a parameter
	result := db.Where("status = ?", "scheduled").Find(&schedules)

	if result.Error != nil {
		return nil, result.Error
	}

	return schedules, nil
}
