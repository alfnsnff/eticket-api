package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type RouteRepository struct {
	DB *gorm.DB
}

func NewRouteRepository(db *gorm.DB) entities.RouteRepositoryInterface {
	return &RouteRepository{DB: db}
}

// Create inserts a new route into the database
func (r *RouteRepository) Create(route *entities.Route) error {
	result := r.DB.Create(route)
	return result.Error
}

// GetAll retrieves all routes from the database
func (r *RouteRepository) GetAll() ([]*entities.Route, error) {
	var Routes []*entities.Route
	result := r.DB.Preload("DepartureHarbor").Preload("ArrivalHarbor").Find(&Routes) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return Routes, nil
}

func (r *RouteRepository) Search(departureHarborID uint, arrivalHarborID uint) (*entities.Route, error) {
	var route entities.Route

	result := r.DB.Where("departure_harbor_id = ? AND arrival_harbor_id = ?", departureHarborID, arrivalHarborID).First(&route)

	if result.Error != nil {
		return nil, result.Error
	}

	return &route, nil
}

// GetByID retrieves a route by its ID
func (r *RouteRepository) GetByID(id uint) (*entities.Route, error) {
	var Route entities.Route
	result := r.DB.Preload("DepartureHarbor").Preload("ArrivalHarbor").First(&Route, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &Route, result.Error
}

// Update modifies an existing route in the database
func (r *RouteRepository) Update(route *entities.Route) error {
	// Uses Gorm's Save method to update the Route
	result := r.DB.Save(route)
	return result.Error
}

// Delete removes a route from the database by its ID
func (r *RouteRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Route{}, id) // Deletes the route by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no route found to delete") // Custom error for non-existent ID
	}
	return nil
}
