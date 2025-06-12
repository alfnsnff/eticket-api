package class

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type ClassUsecase struct {
	Tx              *tx.TxManager
	ClassRepository *repository.ClassRepository
}

func NewClassUsecase(
	tx *tx.TxManager,
	class_repository *repository.ClassRepository,
) *ClassUsecase {
	return &ClassUsecase{
		Tx:              tx,
		ClassRepository: class_repository,
	}
}

func (c *ClassUsecase) CreateClass(ctx context.Context, request *model.WriteClassRequest) error {
	class := mapper.ClassMapper.FromWrite(request)

	if class.ClassName == "" {
		return fmt.Errorf("class name cannot be empty")
	}

	return c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return c.ClassRepository.Create(tx, class)
	})
}

func (c *ClassUsecase) GetAllClasses(ctx context.Context, limit, offset int) ([]*model.ReadClassResponse, int, error) {
	classes := []*entity.Class{}
	var total int64

	err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = c.ClassRepository.Count(tx)
		if err != nil {
			return err
		}
		classes, err = c.ClassRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ClassMapper.ToModels(classes), int(total), nil

}

// GetClassByID retrieves a class by its ID
func (c *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*model.ReadClassResponse, error) {
	class := new(entity.Class)

	err := c.Tx.Execute(ctx, func(tx *gorm.DB) error {
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

	return c.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return c.ClassRepository.Update(tx, class)
	})

}

// DeleteClass deletes a class by its ID
func (c *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {

	return c.Tx.Execute(ctx, func(tx *gorm.DB) error {
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
