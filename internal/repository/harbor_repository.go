package repository

import (
	"errors"
	"eticket-api/internal/domain"

	"gorm.io/gorm"
)

type HarborRepository struct {
	DB *gorm.DB
}

func NewHarborRepository(db *gorm.DB) domain.HarborRepositoryInterface {
	return &HarborRepository{DB: db}
}

// Create inserts a new harbor into the database
func (r *HarborRepository) Create(harbor *domain.Harbor) error {
	result := r.DB.Create(harbor)
	return result.Error
}

// GetAll retrieves all harbor from the database
func (r *HarborRepository) GetAll() ([]*domain.Harbor, error) {
	var harbors []*domain.Harbor
	result := r.DB.Find(&harbors) // Preloads harbor relationship
	// result := r.DB.Find(&harbor)
	if result.Error != nil {
		return nil, result.Error
	}
	return harbors, nil
}

// GetByID retrieves a harbor by its ID
func (r *HarborRepository) GetByID(id uint) (*domain.Harbor, error) {
	var harbor domain.Harbor
	result := r.DB.First(&harbor, id) // Preloads harbor and fetches by ID
	// result := r.DB.First(&harbor, id) // Fetches the harbor by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no harbor is found
	}
	return &harbor, result.Error
}

// Update modifies an existing harbor in the database
func (r *HarborRepository) Update(harbor *domain.Harbor) error {
	// Uses Gorm's Save method to update the harbor
	result := r.DB.Save(harbor)
	return result.Error
}

// Delete removes a harbor from the database by its ID
func (r *HarborRepository) Delete(id uint) error {
	result := r.DB.Delete(&domain.Harbor{}, id) // Deletes the harbor by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no harbor found to delete") // Custom error for non-existent ID
	}
	return nil
}
