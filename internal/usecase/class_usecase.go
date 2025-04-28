package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ClassUsecase struct {
	DB              *gorm.DB
	ClassRepository *repository.ClassRepository
}

func NewClassUsecase(db *gorm.DB, classRepository *repository.ClassRepository) *ClassUsecase {
	return &ClassUsecase{DB: db, ClassRepository: classRepository}
}

// CreateClass validates and creates a new class
func (s *ClassUsecase) CreateClass(ctx context.Context, class *entities.Class) error {
	if class.Name == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ClassRepository.Create(txDB, class)
	})
}

// GetAllClasses retrieves all classes
func (s *ClassUsecase) GetAllClasses(ctx context.Context) ([]*entities.Class, error) {
	// return s.ClassRepository.GetAll()

	var classes []*entities.Class

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		classes, err = s.ClassRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return classes, nil

}

// GetClassByID retrieves a class by its ID
func (s *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*entities.Class, error) {
	// class, err := s.ClassRepository.GetByID(id)

	var class *entities.Class

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		class, err = s.ClassRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}
	return class, nil
}

// UpdateClass updates an existing class
func (s *ClassUsecase) UpdateClass(ctx context.Context, id uint, class *entities.Class) error {
	class.ID = id

	if class.ID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	if class.Name == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ClassRepository.Update(txDB, class)
	})

}

// DeleteClass deletes a class by its ID
func (s *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		class, err := s.ClassRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if class == nil {
			return errors.New("route not found")
		}
		return s.ClassRepository.Delete(txDB, id)
	})

}
