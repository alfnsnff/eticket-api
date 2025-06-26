package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type HarborUsecase struct {
	DB               *gorm.DB
	HarborRepository domain.HarborRepository
}

func NewHarborUsecase(
	db *gorm.DB,
	harborRepository domain.HarborRepository,
) *HarborUsecase {
	return &HarborUsecase{
		DB:               db,
		HarborRepository: harborRepository,
	}
}

func (h *HarborUsecase) CreateHarbor(ctx context.Context, request *model.WriteHarborRequest) error {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	harbor := &domain.Harbor{
		HarborName:    request.HarborName,
		Status:        request.Status,
		HarborAlias:   request.HarborAlias,
		YearOperation: request.YearOperation,
	}

	if err := h.HarborRepository.Insert(tx, harbor); err != nil {
		return fmt.Errorf("failed to create harbor: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (h *HarborUsecase) ListHarbors(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadHarborResponse, int, error) {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := h.HarborRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count harbors: %w", err)
	}

	harbors, err := h.HarborRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get harbors: %w", err)
	}

	responses := make([]*model.ReadHarborResponse, len(harbors))
	for i, harbor := range harbors {
		responses[i] = mapper.HarborToResponse(harbor)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (h *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	harbor, err := h.HarborRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get harbor: %w", err)
	}
	if harbor == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.HarborToResponse(harbor), nil
}

func (h *HarborUsecase) UpdateHarbor(ctx context.Context, request *model.UpdateHarborRequest) error {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing allocation
	harbor, err := h.HarborRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find harbor: %w", err)
	}
	if harbor == nil {
		return errs.ErrNotFound
	}

	harbor.HarborName = request.HarborName
	harbor.Status = request.Status
	harbor.HarborAlias = request.HarborAlias
	harbor.YearOperation = request.YearOperation

	if err := h.HarborRepository.Update(tx, harbor); err != nil {
		return fmt.Errorf("failed to update harbor: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (h *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	harbor, err := h.HarborRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get harbor: %w", err)
	}
	if harbor == nil {
		return errs.ErrNotFound
	}

	if err := h.HarborRepository.Delete(tx, harbor); err != nil {
		return fmt.Errorf("failed to delete harbor: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
