package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type HarborRepository struct {
	Repository[entities.Harbor]
}

func NewHarborRepository() *HarborRepository {
	return &HarborRepository{}
}

// GetAll retrieves all harbor from the database
func (r *HarborRepository) GetAll(db *gorm.DB) ([]*entities.Harbor, error) {
	var harbors []*entities.Harbor
	result := db.Find(&harbors) // Preloads harbor relationship
	// result := r.DB.Find(&harbor)
	if result.Error != nil {
		return nil, result.Error
	}
	return harbors, nil
}

// GetByID retrieves a harbor by its ID
func (r *HarborRepository) GetByID(db *gorm.DB, id uint) (*entities.Harbor, error) {
	var harbor entities.Harbor
	result := db.First(&harbor, id) // Preloads harbor and fetches by ID
	// result := r.DB.First(&harbor, id) // Fetches the harbor by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no harbor is found
	}
	return &harbor, result.Error
}
