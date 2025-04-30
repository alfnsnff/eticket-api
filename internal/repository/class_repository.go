package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ClassRepository struct {
	Repository[entities.Class]
}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

// GetAll retrieves all classes from the database
func (r *ClassRepository) GetAll(db *gorm.DB) ([]*entities.Class, error) {
	var classes []*entities.Class
	result := db.Find(&classes) // Preloads Class relationship
	// result := r.DB.Find(&classes)
	if result.Error != nil {
		return nil, result.Error
	}
	return classes, nil
}

// GetByID retrieves a class by its ID
func (r *ClassRepository) GetByID(db *gorm.DB, id uint) (*entities.Class, error) {
	var class entities.Class
	result := db.First(&class, id) // Preloads Class and fetches by ID
	// result := r.DB.First(&class, id) // Fetches the class by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no class is found
	}
	return &class, result.Error
}
