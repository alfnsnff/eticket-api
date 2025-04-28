package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ShipRepository struct {
	DB *gorm.DB
}

func NewShipRepository() *ShipRepository {
	return &ShipRepository{}
}

// Create inserts a new ship into the database
func (r *ShipRepository) Create(db *gorm.DB, ship *entities.Ship) error {
	result := db.Create(ship)
	return result.Error
}

// GetAll retrieves all ships from the database
func (r *ShipRepository) GetAll(db *gorm.DB) ([]*entities.Ship, error) {
	var ships []*entities.Ship
	result := db.Preload("ShipClasses").Preload("ShipClasses.Class").Preload("ShipClasses.Ship").Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

// GetByID retrieves a ship by its ID
func (r *ShipRepository) GetByID(db *gorm.DB, id uint) (*entities.Ship, error) {
	var ship entities.Ship
	result := db.Preload("ShipClasses").Preload("ShipClasses.Class").Preload("ShipClasses.Ship").First(&ship, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &ship, result.Error
}

// Update modifies an existing ship in the database
func (r *ShipRepository) Update(db *gorm.DB, ship *entities.Ship) error {
	// Uses Gorm's Save method to update the ship
	result := db.Save(ship)
	return result.Error
}

// Delete removes a ship from the database by its ID
func (r *ShipRepository) Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&entities.Ship{}, id) // Deletes the ship by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ship found to delete") // Custom error for non-existent ID
	}
	return nil
}
