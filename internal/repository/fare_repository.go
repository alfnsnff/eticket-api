package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type FareRepository struct {
	Repository[entities.Fare]
}

func NewFareRepository() *FareRepository {
	return &FareRepository{}
}

// GetAll retrieves all routes from the database
func (p *FareRepository) GetAll(db *gorm.DB) ([]*entities.Fare, error) {
	var fares []*entities.Fare
	result := db.Find(&fares) // Corrected Preload
	if result.Error != nil {
		return nil, result.Error
	}
	return fares, nil
}

// GetByID retrieves a route by its ID
func (p *FareRepository) GetByID(db *gorm.DB, id uint) (*entities.Fare, error) {
	var fare entities.Fare
	result := db.First(&fare, id) // Fetches the route by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return &fare, result.Error
}

// GetByID retrieves a route by its ID
func (p *FareRepository) GetByIDs(db *gorm.DB, ids []uint) ([]*entities.Fare, error) {
	var fares []*entities.Fare

	result := db.Where("id IN ?", ids).Preload("Manifest").
		Preload("Manifest.Class").Find(&fares) // Fetches the route by ID

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return fares, result.Error
}

// FareRepository.go
func (r *FareRepository) GetByRouteID(db *gorm.DB, routeID uint) ([]*entities.Fare, error) {
	var fares []*entities.Fare
	result := db.
		Preload("Manifest").
		Preload("Manifest.Class"). // for ClassName
		Where("route_id = ?", routeID).
		Find(&fares)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no Route is found
	}
	return fares, result.Error
}
