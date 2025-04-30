package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ClassUsecase struct {
	DB              *gorm.DB
	ClassRepository *repository.ClassRepository
}

func NewClassUsecase(db *gorm.DB, class_repository *repository.ClassRepository) *ClassUsecase {
	return &ClassUsecase{DB: db, ClassRepository: class_repository}
}

// CreateClass validates and creates a new class
func (s *ClassUsecase) CreateClass(ctx context.Context, request *model.WriteClassRequest) error {
	class := mapper.ToClassEntity(request)

	if class.Name == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ClassRepository.Create(tx, class)
	})
}

// GetAllClasses retrieves all classes
func (s *ClassUsecase) GetAllClasses(ctx context.Context) ([]*model.ReadClassResponse, error) {
	// return s.ClassRepository.GetAll()

	classes := []*entities.Class{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		classes, err = s.ClassRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ToClassesModel(classes), nil

}

// GetClassByID retrieves a class by its ID
func (s *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*model.ReadClassResponse, error) {
	// class, err := s.ClassRepository.GetByID(id)

	class := new(entities.Class)
	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		class, err = s.ClassRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}
	return mapper.ToClassModel(class), nil
}

// UpdateClass updates an existing class
func (s *ClassUsecase) UpdateClass(ctx context.Context, id uint, request *model.WriteClassRequest) error {
	class := mapper.ToClassEntity(request)
	class.ID = id

	if class.ID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	if class.Name == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ClassRepository.Update(tx, class)
	})

}

// DeleteClass deletes a class by its ID
func (s *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		class, err := s.ClassRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if class == nil {
			return errors.New("route not found")
		}
		return s.ClassRepository.Delete(tx, class)
	})

}
