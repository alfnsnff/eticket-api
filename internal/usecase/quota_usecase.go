package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type QuotaUsecase struct {
	Transactor      *transact.Transactor
	QuotaRepository domain.QuotaRepository
}

func NewQuotaUsecase(

	transactor *transact.Transactor,
	Quota_repository domain.QuotaRepository,
) *QuotaUsecase {
	return &QuotaUsecase{

		Transactor:      transactor,
		QuotaRepository: Quota_repository,
	}
}

func (uc *QuotaUsecase) CreateQuota(ctx context.Context, e *domain.Quota) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		quota := &domain.Quota{
			ScheduleID: e.ScheduleID,
			ClassID:    e.ClassID,
			Quota:      e.Capacity,
			Capacity:   e.Capacity,
			Price:      e.Price,
		}
		if err := uc.QuotaRepository.Insert(ctx, tx, quota); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create quota: %w", err)
		}
		return nil
	})
}

func (uc *QuotaUsecase) CreateQuotaBulk(ctx context.Context, es []*domain.Quota) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		quotas := make([]*domain.Quota, len(es))
		for i, e := range es {
			quotas[i] = &domain.Quota{
				ScheduleID: e.ScheduleID,
				ClassID:    e.ClassID,
				Price:      e.Price,
				Quota:      e.Capacity,
				Capacity:   e.Capacity,
			}
		}

		if err := uc.QuotaRepository.InsertBulk(ctx, tx, quotas); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create quotas in bulk: %w", err)
		}

		return nil
	})
}

func (uc *QuotaUsecase) ListQuotas(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Quota, int, error) {
	var err error
	var total int64
	var quotas []*domain.Quota
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.QuotaRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count Quotas: %w", err)
		}
		quotas, err = uc.QuotaRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all Quotas: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list Quotas: %w", err)
	}

	return quotas, int(total), nil
}

func (uc *QuotaUsecase) GetQuotaByID(ctx context.Context, id uint) (*domain.Quota, error) {
	var err error
	var quota *domain.Quota
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		quota, err = uc.QuotaRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get Quota: %w", err)
		}

		if quota == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get Quota by ID: %w", err)
	}

	return quota, nil
}

func (uc *QuotaUsecase) UpdateQuota(ctx context.Context, e *domain.Quota) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		quota, err := uc.QuotaRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find quota: %w", err)
		}
		if quota == nil {
			return errs.ErrNotFound
		}

		quota.ScheduleID = e.ScheduleID
		quota.ClassID = e.ClassID
		quota.Quota = e.Capacity
		quota.Capacity = e.Capacity
		quota.Price = e.Price

		if err := uc.QuotaRepository.Update(ctx, tx, quota); err != nil {
			return fmt.Errorf("failed to update quota: %w", err)
		}
		return nil
	})
}

func (uc *QuotaUsecase) DeleteQuota(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		quota, err := uc.QuotaRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get quota: %w", err)
		}
		if quota == nil {
			return errs.ErrNotFound
		}

		if err := uc.QuotaRepository.Delete(ctx, tx, quota); err != nil {
			return fmt.Errorf("failed to delete quota: %w", err)
		}
		return nil
	})
}
