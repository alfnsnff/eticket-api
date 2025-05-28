package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
)

type RouteRepository struct {
	Repository[entity.Route]
}

func NewRouteRepository() *RouteRepository {
	return &RouteRepository{}
}

func (rr *RouteRepository) Count(db *gorm.DB) (int64, error) {
	routes := []*entity.Route{}
	var total int64
	result := db.Find(&routes).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (rr *RouteRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Route, error) {
	routes := []*entity.Route{}
	result := db.Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		Limit(limit).Offset(offset).Find(&routes)
	if result.Error != nil {
		return nil, result.Error
	}
	return routes, nil
}

func (rr *RouteRepository) GetByID(db *gorm.DB, id uint) (*entity.Route, error) {
	route := new(entity.Route)
	result := db.Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		First(&route, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return route, result.Error
}
