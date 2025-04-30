package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ManifestRepository struct {
	Repository[entities.Manifest]
}

func NewManifestRepository() *ManifestRepository {
	return &ManifestRepository{}
}

// GetAll retrieves all ships from the database
func (sc *ManifestRepository) GetAll(db *gorm.DB) ([]*entities.Manifest, error) {
	var manifests []*entities.Manifest
	result := db.Preload("Class").Preload("Ship").Find(&manifests)
	if result.Error != nil {
		return nil, result.Error
	}
	return manifests, nil
}

func (sc *ManifestRepository) GetByShipAndClass(db *gorm.DB, shipID, classID uint) (*entities.Manifest, error) {
	var manifest entities.Manifest
	err := db.Where("ship_id = ? AND class_id = ?", shipID, classID).First(&manifest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &manifest, err
}

// GetByID retrieves a ship by its ID
func (sc *ManifestRepository) GetByID(db *gorm.DB, id uint) (*entities.Manifest, error) {
	var manifest entities.Manifest
	result := db.Preload("Class").Preload("Ship").First(&manifest, id) // Fetches the ship by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no ship is found
	}
	return &manifest, result.Error
}

func (sc *ManifestRepository) GetByIDs(db *gorm.DB, ids []uint) ([]*entities.Manifest, error) {
	var manifests []*entities.Manifest
	err := db.Where("id IN ?", ids).Preload("Class").Preload("Ship").Find(&manifests).Error
	if err != nil {
		return nil, err
	}
	return manifests, nil
}

func (sc *ManifestRepository) GetByShipID(db *gorm.DB, shipId uint) ([]*entities.Manifest, error) {
	var manifests []*entities.Manifest
	result := db.Where("ship_id = ?", shipId).Preload("Class").Preload("Ship").Find(&manifests)

	if result.Error != nil {
		return nil, result.Error
	}

	return manifests, nil
}
