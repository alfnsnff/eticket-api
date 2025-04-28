package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type RouteRepository struct {
	DB *gorm.DB
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

// Create inserts a new route into the database
func (r *RouteRepository) Create(db *gorm.DB, route *entities.Route) error {
	result := db.Create(route)
	return result.Error
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

// Update modifies an existing route in the database
func (r *RouteRepository) Update(db *gorm.DB, route *entities.Route) error {
	// Uses Gorm's Save method to update the Route
	result := db.Save(route)
	return result.Error
}

// Delete removes a route from the database by its ID
func (r *RouteRepository) Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&entities.Route{}, id) // Deletes the route by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no route found to delete") // Custom error for non-existent ID
	}
	return nil
}
