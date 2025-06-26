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

type ClaimItemUsecase struct {
	DB                  *gorm.DB
	ClaimItemRepository domain.ClaimItemRepository
}

func NewClaimItemUsecase(
	db *gorm.DB,
	claim_item_repository domain.ClaimItemRepository,
) *ClaimItemUsecase {
	return &ClaimItemUsecase{
		DB:                  db,
		ClaimItemRepository: claim_item_repository,
	}
}

func (c *ClaimItemUsecase) CreateClaimItem(ctx context.Context, request *model.WriteClaimItemRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	claimItem := &domain.ClaimItem{
		ClaimSessionID: request.ClaimSessionID,
		ClassID:        request.ClassID,
		Quantity:       request.Quantity,
	}

	if err := c.ClaimItemRepository.Insert(tx, claimItem); err != nil {
		return fmt.Errorf("failed to create claim item: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *ClaimItemUsecase) ListClaimItemes(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClaimItemResponse, int, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := c.ClaimItemRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ClaimItemes: %w", err)
	}

	ClaimItemes, err := c.ClaimItemRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all ClaimItemes: %w", err)
	}

	responses := make([]*model.ReadClaimItemResponse, len(ClaimItemes))
	for i, ClaimItem := range ClaimItemes {
		responses[i] = mapper.ClaimItemToResponse(ClaimItem)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (c *ClaimItemUsecase) GetClaimItemByID(ctx context.Context, id uint) (*model.ReadClaimItemResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	ClaimItem, err := c.ClaimItemRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim item: %w", err)
	}
	if ClaimItem == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ClaimItemToResponse(ClaimItem), nil
}

func (c *ClaimItemUsecase) UpdateClaimItem(ctx context.Context, request *model.UpdateClaimItemRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	claimItem, err := c.ClaimItemRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find claim item: %w", err)
	}
	if claimItem == nil {
		return errs.ErrNotFound
	}

	claimItem.ClaimSessionID = uint(request.Quantity)
	claimItem.ClassID = request.ClassID
	claimItem.Quantity = request.Quantity

	if err := c.ClaimItemRepository.Update(tx, claimItem); err != nil {
		return fmt.Errorf("failed to update claim item: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (c *ClaimItemUsecase) DeleteClaimItem(ctx context.Context, id uint) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	ClaimItem, err := c.ClaimItemRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get claim item: %w", err)
	}
	if ClaimItem == nil {
		return errs.ErrNotFound
	}

	if err := c.ClaimItemRepository.Delete(tx, ClaimItem); err != nil {
		return fmt.Errorf("failed to delete claim item: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
