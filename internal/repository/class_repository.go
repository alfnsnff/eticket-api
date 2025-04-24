package repository

import (
	"errors"
	"eticket-api/internal/domain/entities"

	"gorm.io/gorm"
)

type ClassRepository struct {
	DB *gorm.DB
}

func NewClassRepository(db *gorm.DB) entities.ClassRepositoryInterface {
	return &ClassRepository{DB: db}
}

// Create inserts a new class into the database
func (r *ClassRepository) Create(class *entities.Class) error {
	result := r.DB.Create(class)
	return result.Error
}

// GetAll retrieves all classes from the database
func (r *ClassRepository) GetAll() ([]*entities.Class, error) {
	var classes []*entities.Class
	result := r.DB.Find(&classes) // Preloads Class relationship
	// result := r.DB.Find(&classes)
	if result.Error != nil {
		return nil, result.Error
	}
	return classes, nil
}

// GetByID retrieves a class by its ID
func (r *ClassRepository) GetByID(id uint) (*entities.Class, error) {
	var class entities.Class
	result := r.DB.First(&class, id) // Preloads Class and fetches by ID
	// result := r.DB.First(&class, id) // Fetches the class by ID
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Returns nil if no class is found
	}
	return &class, result.Error
}

// Update modifies an existing class in the database
func (r *ClassRepository) Update(class *entities.Class) error {
	// Uses Gorm's Save method to update the class
	result := r.DB.Save(class)
	return result.Error
}

// Delete removes a class from the database by its ID
func (r *ClassRepository) Delete(id uint) error {
	result := r.DB.Delete(&entities.Class{}, id) // Deletes the class by ID
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no class found to delete") // Custom error for non-existent ID
	}
	return nil
}
