package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type RouteRepository struct {
	Repository[entities.Route]
}

func NewRouteRepository() *RouteRepository {
	return &RouteRepository{}
}

// GetAll retrieves all routes from the database
func (r *RouteRepository) GetAll(db *gorm.DB) ([]*entities.Route, error) {
	var Routes []*entities.Route
	result := db.Preload("DepartureHarbor").Preload("ArrivalHarbor").Find(&Routes) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return Routes, nil
}

func (r *RouteRepository) Search(db *gorm.DB, departureHarborID uint, arrivalHarborID uint) (*entities.Route, error) {
	var route entities.Route

	result := db.Where("departure_harbor_id = ? AND arrival_harbor_id = ?", departureHarborID, arrivalHarborID).First(&route)

	if result.Error != nil {
		return nil, result.Error
	}

	return &route, nil
}

// GetByID retrieves a route by its ID
func (r *RouteRepository) GetByID(db *gorm.DB, id uint) (*entities.Route, error) {
	var Route entities.Route
	result := db.Preload("DepartureHarbor").Preload("ArrivalHarbor").First(&Route, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &Route, result.Error
}
