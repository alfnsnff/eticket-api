package harbor

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type HarborUsecase struct {
	DB               *gorm.DB
	HarborRepository HarborRepository
}

func NewHarborUsecase(
	db *gorm.DB,
	harborRepository HarborRepository,
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
		} else {
			tx.Rollback()
		}
	}()

	harbor := &entity.Harbor{
		HarborName:    request.HarborName,
		Status:        request.Status,
		HarborAlias:   request.HarborAlias,
		YearOperation: request.YearOperation,
	}

	if err := h.HarborRepository.Create(tx, harbor); err != nil {
		return fmt.Errorf("failed to create harbor: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (h *HarborUsecase) GetAllHarbors(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadHarborResponse, int, error) {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := h.HarborRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count harbors: %w", err)
	}

	harbors, err := h.HarborRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get harbors: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadHarborResponses(harbors), int(total), nil
}

func (h *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	harbor, err := h.HarborRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get harbor: %w", err)
	}
	if harbor == nil {
		return nil, errors.New("harbor not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return ToReadHarborResponse(harbor), nil
}

func (h *HarborUsecase) UpdateHarbor(ctx context.Context, request *model.UpdateHarborRequest) error {
	tx := h.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	harbor, err := h.HarborRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find harbor: %w", err)
	}
	if harbor == nil {
		return errors.New("harbor not found")
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
		} else {
			tx.Rollback()
		}
	}()

	harbor, err := h.HarborRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get harbor: %w", err)
	}
	if harbor == nil {
		return errors.New("harbor not found")
	}

	if err := h.HarborRepository.Delete(tx, harbor); err != nil {
		return fmt.Errorf("failed to delete harbor: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
