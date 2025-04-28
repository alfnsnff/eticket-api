package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ShipClassRepository struct {
	DB *gorm.DB
}

func NewShipClassRepository() *ShipClassRepository {
	return &ShipClassRepository{}
}

// Create inserts a new ship into the database
func (sc *ShipClassRepository) Create(db *gorm.DB, shipClass *entities.ShipClass) error {
	result := db.Create(shipClass)
	return result.Error
}

// GetAll retrieves all ships from the database
func (sc *ShipClassRepository) GetAll(db *gorm.DB) ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	result := db.Preload("Class").Preload("Ship").Find(&shipClasses)
	if result.Error != nil {
		return nil, result.Error
	}
	return shipClasses, nil
}

func (sc *ShipClassRepository) GetByShipAndClass(db *gorm.DB, shipID, classID uint) (*entities.ShipClass, error) {
	var shipClass entities.ShipClass
	err := db.Where("ship_id = ? AND class_id = ?", shipID, classID).First(&shipClass).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &shipClass, err
}

// GetByID retrieves a ship by its ID
func (sc *ShipClassRepository) GetByID(db *gorm.DB, id uint) (*entities.ShipClass, error) {
	var shipClass entities.ShipClass
	result := db.Preload("Class").Preload("Ship").First(&shipClass, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &shipClass, result.Error
}

func (sc *ShipClassRepository) GetByIDs(db *gorm.DB, ids []uint) ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	err := db.Where("id IN ?", ids).Preload("Class").Preload("Ship").Find(&shipClasses).Error
	if err != nil {
		return nil, err
	}
	return shipClasses, nil
}

func (sc *ShipClassRepository) GetByShipID(db *gorm.DB, shipId uint) ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	result := db.Where("ship_id = ?", shipId).Preload("Class").Preload("Ship").Find(&shipClasses)

	if result.Error != nil {
		return nil, result.Error
	}

	return shipClasses, nil
}

// Update modifies an existing ship in the database
func (sc *ShipClassRepository) Update(db *gorm.DB, shipClass *entities.ShipClass) error {
	// Uses Gorm's Save method to update the ship
	result := db.Save(shipClass)
	return result.Error
}

// Delete removes a ship from the database by its ID
func (sc *ShipClassRepository) Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&entities.ShipClass{}, id) // Deletes the ship by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ship found to delete") // Custom error for non-existent ID
	}
	return nil
}
