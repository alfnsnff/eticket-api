package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RouteRepository struct{}

func NewRouteRepository() *RouteRepository {
	return &RouteRepository{}
}

func (ar *RouteRepository) Create(db *gorm.DB, route *domain.Route) error {
	result := db.Create(route)
	return result.Error
}

func (ar *RouteRepository) Update(db *gorm.DB, route *domain.Route) error {
	result := db.Save(route)
	return result.Error
}

func (ar *RouteRepository) Delete(db *gorm.DB, route *domain.Route) error {
	result := db.Select(clause.Associations).Delete(route)
	return result.Error
}

func (rr *RouteRepository) Count(db *gorm.DB) (int64, error) {
	routes := []*domain.Route{}
	var total int64
	result := db.Find(&routes).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (rr *RouteRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.Route, error) {
	routes := []*domain.Route{}

	query := db.Preload("DepartureHarbor").Preload("ArrivalHarbor")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("departure_harbor ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&routes).Error
	return routes, err
}

func (rr *RouteRepository) GetByID(db *gorm.DB, id uint) (*domain.Route, error) {
	route := new(domain.Route)
	result := db.Preload("DepartureHarbor").
		Preload("ArrivalHarbor").
		First(&route, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return route, result.Error
}
