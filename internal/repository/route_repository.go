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

func (rr *RouteRepository) GetAll(db *gorm.DB) ([]*entity.Route, error) {
	routes := []*entity.Route{}
	result := db.Find(&routes)
	if result.Error != nil {
		return nil, result.Error
	}
	return routes, nil
}

func (rr *RouteRepository) GetByID(db *gorm.DB, id uint) (*entity.Route, error) {
	route := new(entity.Route)
	result := db.First(&route, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return route, result.Error
}
