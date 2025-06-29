package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type ClassUsecase struct {
	Transactor      *transact.Transactor
	ClassRepository domain.ClassRepository
}

func NewClassUsecase(
	transactor *transact.Transactor,
	class_repository domain.ClassRepository,
) *ClassUsecase {
	return &ClassUsecase{
		Transactor:      transactor,
		ClassRepository: class_repository,
	}
}

func (uc *ClassUsecase) CreateClass(ctx context.Context, e *domain.Class) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		class := &domain.Class{
			ClassName:  e.ClassName,
			Type:       e.Type,
			ClassAlias: e.ClassAlias,
		}

		if err := uc.ClassRepository.Insert(ctx, tx, class); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create class: %w", err)
		}
		return nil
	})
}

func (uc *ClassUsecase) ListClasses(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Class, int, error) {

	var err error
	var total int64
	var classes []*domain.Class
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.ClassRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count classes: %w", err)
		}
		classes, err = uc.ClassRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all classes: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list classes: %w", err)
	}

	return classes, int(total), nil
}

func (uc *ClassUsecase) GetClassByID(ctx context.Context, id uint) (*domain.Class, error) {
	var err error
	var class *domain.Class
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		class, err = uc.ClassRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get class: %w", err)
		}
		if class == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get class by ID: %w", err)
	}
	return class, nil
}

func (uc *ClassUsecase) UpdateClass(ctx context.Context, e *domain.Class) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		class, err := uc.ClassRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find class: %w", err)
		}
		if class == nil {
			return errs.ErrNotFound
		}

		class.ClassName = e.ClassName
		class.Type = e.Type
		class.ClassAlias = e.ClassAlias

		if err := uc.ClassRepository.Update(ctx, tx, class); err != nil {
			return fmt.Errorf("failed to update class: %w", err)
		}

		return nil
	})
}

func (uc *ClassUsecase) DeleteClass(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {

		class, err := uc.ClassRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get class: %w", err)
		}
		if class == nil {
			return errs.ErrNotFound
		}

		if err := uc.ClassRepository.Delete(ctx, tx, class); err != nil {
			return fmt.Errorf("failed to delete class: %w", err)
		}

		return nil
	})
}
