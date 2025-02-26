package repository

import (
	"errors"
	"eticket-api/internal/domain"

	"gorm.io/gorm"
)

type RouteRepository struct {
	DB *gorm.DB
}

func NewRouteRepository(db *gorm.DB) domain.RouteRepositoryInterface {
	return &RouteRepository{DB: db}
}

// Create inserts a new route into the database
func (r *RouteRepository) Create(route *domain.Route) error {
	result := r.DB.Create(route)
	return result.Error
}

// GetAll retrieves all routes from the database
func (r *RouteRepository) GetAll() ([]*domain.Route, error) {
	var Routes []*domain.Route
	result := r.DB.Preload("DepartureHarbor").Preload("ArrivalHarbor").Find(&Routes) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return Routes, nil
}

// GetByID retrieves a route by its ID
func (r *RouteRepository) GetByID(id uint) (*domain.Route, error) {
	var Route domain.Route
	result := r.DB.Preload("DepartureHarbor").Preload("ArrivalHarbor").First(&Route, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &Route, result.Error
}

// Update modifies an existing route in the database
func (r *RouteRepository) Update(route *domain.Route) error {
	// Uses Gorm's Save method to update the Route
	result := r.DB.Save(route)
	return result.Error
}

// Delete removes a route from the database by its ID
func (r *RouteRepository) Delete(id uint) error {
	result := r.DB.Delete(&domain.Route{}, id) // Deletes the route by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no route found to delete") // Custom error for non-existent ID
	}
	return nil
}
