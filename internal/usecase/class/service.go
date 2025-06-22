package class

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type ClassUsecase struct {
	DB              *gorm.DB
	ClassRepository ClassRepository
}

func NewClassUsecase(
	db *gorm.DB,
	class_repository ClassRepository,
) *ClassUsecase {
	return &ClassUsecase{
		DB:              db,
		ClassRepository: class_repository,
	}
}

func (c *ClassUsecase) CreateClass(ctx context.Context, request *model.WriteClassRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	class := &entity.Class{
		ClassName:  request.ClassName,
		Type:       request.Type,
		ClassAlias: request.ClassAlias,
	}

	if err := c.ClassRepository.Create(tx, class); err != nil {
		return fmt.Errorf("failed to create class: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *ClassUsecase) GetAllClasses(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClassResponse, int, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := c.ClassRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count classes: %w", err)
	}

	classes, err := c.ClassRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all classes: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadClassResponses(classes), int(total), nil
}

func (c *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*model.ReadClassResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	class, err := c.ClassRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get class: %w", err)
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadClassResponse(class), nil
}

func (c *ClassUsecase) UpdateClass(ctx context.Context, request *model.UpdateClassRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	class, err := c.ClassRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find class: %w", err)
	}
	if class == nil {
		return errors.New("class not found")
	}

	class.ClassName = request.ClassName
	class.Type = request.Type
	class.ClassAlias = request.ClassAlias

	if err := c.ClassRepository.Update(tx, class); err != nil {
		return fmt.Errorf("failed to update class: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	class, err := c.ClassRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get class: %w", err)
	}
	if class == nil {
		return errors.New("class not found")
	}

	if err := c.ClassRepository.Delete(tx, class); err != nil {
		return fmt.Errorf("failed to delete class: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
