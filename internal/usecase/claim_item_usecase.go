package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"eticket-api/pkg/gotann"
	"fmt"
)

type ClaimItemUsecase struct {
	Transactor          transact.Transactor
	ClaimItemRepository domain.ClaimItemRepository
}

func NewClaimItemUsecase(
	transactor transact.Transactor,
	claim_item_repository domain.ClaimItemRepository,
) *ClaimItemUsecase {
	return &ClaimItemUsecase{
		Transactor:          transactor,
		ClaimItemRepository: claim_item_repository,
	}
}

func (uc *ClaimItemUsecase) CreateClaimItem(ctx context.Context, e *domain.ClaimItem) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimItem := &domain.ClaimItem{
			ClaimSessionID: e.ClaimSessionID,
			ClassID:        e.ClassID,
			Quantity:       e.Quantity,
		}
		if err := uc.ClaimItemRepository.Insert(ctx, tx, claimItem); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create claim item: %w", err)
		}
		return nil
	})
}

func (uc *ClaimItemUsecase) ListClaimItemes(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadClaimItemResponse, int, error) {

	var err error
	var total int64
	var claimItems []*domain.ClaimItem
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.ClaimItemRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count ClaimItemes: %w", err)
		}

		claimItems, err = uc.ClaimItemRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all ClaimItemes: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list claim items: %w", err)
	}

	responses := make([]*model.ReadClaimItemResponse, len(claimItems))
	for i, ClaimItem := range claimItems {
		responses[i] = mapper.ClaimItemToResponse(ClaimItem)
	}

	return responses, int(total), nil
}

func (uc *ClaimItemUsecase) GetClaimItemByID(ctx context.Context, id uint) (*model.ReadClaimItemResponse, error) {
	var err error
	var claimItem *domain.ClaimItem
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimItem, err = uc.ClaimItemRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get claim item: %w", err)
		}
		if claimItem == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get claim item by ID: %w", err)
	}
	return mapper.ClaimItemToResponse(claimItem), nil
}

func (uc *ClaimItemUsecase) UpdateClaimItem(ctx context.Context, e *domain.ClaimItem) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimItem, err := uc.ClaimItemRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find claim item: %w", err)
		}
		if claimItem == nil {
			return errs.ErrNotFound
		}

		claimItem.ClaimSessionID = e.ClaimSessionID
		claimItem.ClassID = e.ClassID
		claimItem.Quantity = e.Quantity

		if err := uc.ClaimItemRepository.Update(ctx, tx, claimItem); err != nil {
			return fmt.Errorf("failed to update claim item: %w", err)
		}
		return nil
	})
}

func (uc *ClaimItemUsecase) DeleteClaimItem(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claimItem, err := uc.ClaimItemRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get claim item: %w", err)
		}
		if claimItem == nil {
			return errs.ErrNotFound
		}

		if err := uc.ClaimItemRepository.Delete(ctx, tx, claimItem); err != nil {
			return fmt.Errorf("failed to delete claim item: %w", err)
		}
		return nil
	})
}
