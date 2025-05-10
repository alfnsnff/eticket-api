package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
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

func (c *ClassUsecase) CreateClass(ctx context.Context, request *model.WriteClassRequest) error {
	class := mapper.ClassMapper.FromWrite(request)

	if class.ClassName == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, c.DB, func(tx *gorm.DB) error {
		return c.ClassRepository.Create(tx, class)
	})
}

func (c *ClassUsecase) GetAllClasses(ctx context.Context) ([]*model.ReadClassResponse, error) {
	classes := []*entity.Class{}

	err := tx.Execute(ctx, c.DB, func(tx *gorm.DB) error {
		var err error
		classes, err = c.ClassRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ClassMapper.ToModels(classes), nil

}

// GetClassByID retrieves a class by its ID
func (c *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*model.ReadClassResponse, error) {
	class := new(entity.Class)

	err := tx.Execute(ctx, c.DB, func(tx *gorm.DB) error {
		var err error
		class, err = c.ClassRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if class == nil {
		return nil, errors.New("class not found")
	}

	return mapper.ClassMapper.ToModel(class), nil
}

func (c *ClassUsecase) UpdateClass(ctx context.Context, id uint, request *model.UpdateClassRequest) error {
	class := mapper.ClassMapper.FromUpdate(request)
	class.ID = id

	if class.ID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	if class.ClassName == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return tx.Execute(ctx, c.DB, func(tx *gorm.DB) error {
		return c.ClassRepository.Update(tx, class)
	})

}

// DeleteClass deletes a class by its ID
func (c *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {

	return tx.Execute(ctx, c.DB, func(tx *gorm.DB) error {
		class, err := c.ClassRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if class == nil {
			return errors.New("route not found")
		}
		return c.ClassRepository.Delete(tx, class)
	})

}
