package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ShipRepository struct {
	Repository[entities.Ship]
}

func NewShipRepository() *ShipRepository {
	return &ShipRepository{}
}

// GetAll retrieves all ships from the database
func (r *ShipRepository) GetAll(db *gorm.DB) ([]*entities.Ship, error) {
	var ships []*entities.Ship
	result := db.Preload("Manifests").Preload("Manifests.Class").Preload("Manifests.Ship").Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

// GetByID retrieves a ship by its ID
func (r *ShipRepository) GetByID(db *gorm.DB, id uint) (*entities.Ship, error) {
	var ship entities.Ship
	result := db.Preload("Manifests").Preload("Manifests.Class").Preload("Manifests.Ship").First(&ship, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &ship, result.Error
}
