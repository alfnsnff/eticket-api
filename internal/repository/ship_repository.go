package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ShipRepository struct {
	DB *gorm.DB
}

func NewShipRepository(db *gorm.DB) entities.ShipRepositoryInterface {
	return &ShipRepository{DB: db}
}

// Create inserts a new ship into the database
func (r *ShipRepository) Create(ship *entities.Ship) error {
	result := r.DB.Create(ship)
	return result.Error
}

// GetAll retrieves all ships from the database
func (r *ShipRepository) GetAll() ([]*entities.Ship, error) {
	var ships []*entities.Ship
	result := r.DB.Find(&ships)
	if result.Error != nil {
		return nil, result.Error
	}
	return ships, nil
}

// GetByID retrieves a ship by its ID
func (r *ShipRepository) GetByID(id uint) (*entities.Ship, error) {
	var ship entities.Ship
	result := r.DB.First(&ship, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &ship, result.Error
}

// Update modifies an existing ship in the database
func (r *ShipRepository) Update(ship *entities.Ship) error {
	// Uses Gorm's Save method to update the ship
	result := r.DB.Save(ship)
	return result.Error
}

// Delete removes a ship from the database by its ID
func (r *ShipRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Ship{}, id) // Deletes the ship by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ship found to delete") // Custom error for non-existent ID
	}
	return nil
}
