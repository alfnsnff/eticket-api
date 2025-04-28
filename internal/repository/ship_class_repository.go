package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ShipClassRepository struct {
	DB *gorm.DB
}

func NewShipClassRepository(db *gorm.DB) *ShipClassRepository {
	return &ShipClassRepository{DB: db}
}

// Create inserts a new ship into the database
func (sc *ShipClassRepository) Create(shipClass *entities.ShipClass) error {
	result := sc.DB.Create(shipClass)
	return result.Error
}

// GetAll retrieves all ships from the database
func (sc *ShipClassRepository) GetAll() ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	result := sc.DB.Preload("Class").Preload("Ship").Find(&shipClasses)
	if result.Error != nil {
		return nil, result.Error
	}
	return shipClasses, nil
}

func (sc *ShipClassRepository) GetByShipAndClass(shipID, classID uint) (*entities.ShipClass, error) {
	var shipClass entities.ShipClass
	err := sc.DB.Where("ship_id = ? AND class_id = ?", shipID, classID).First(&shipClass).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &shipClass, err
}

// GetByID retrieves a ship by its ID
func (sc *ShipClassRepository) GetByID(id uint) (*entities.ShipClass, error) {
	var shipClass entities.ShipClass
	result := sc.DB.Preload("Class").Preload("Ship").First(&shipClass, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &shipClass, result.Error
}

func (sc *ShipClassRepository) GetByIDs(ids []uint) ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	err := sc.DB.Where("id IN ?", ids).Preload("Class").Preload("Ship").Find(&shipClasses).Error
	if err != nil {
		return nil, err
	}
	return shipClasses, nil
}

func (sc *ShipClassRepository) GetByShipID(shipId uint) ([]*entities.ShipClass, error) {
	var shipClasses []*entities.ShipClass
	result := sc.DB.Where("ship_id = ?", shipId).Preload("Class").Preload("Ship").Find(&shipClasses)

	if result.Error != nil {
		return nil, result.Error
	}

	return shipClasses, nil
}

// Update modifies an existing ship in the database
func (sc *ShipClassRepository) Update(shipClass *entities.ShipClass) error {
	// Uses Gorm's Save method to update the ship
	result := sc.DB.Save(shipClass)
	return result.Error
}

// Delete removes a ship from the database by its ID
func (sc *ShipClassRepository) Delete(id uint) error {
	result := sc.DB.Delete(&entities.ShipClass{}, id) // Deletes the ship by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no ship found to delete") // Custom error for non-existent ID
	}
	return nil
}
