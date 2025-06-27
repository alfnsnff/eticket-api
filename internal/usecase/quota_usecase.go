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

type QuotaUsecase struct {
	DB              *gorm.DB // Assuming you have a DB field for the transaction manager
	QuotaRepository domain.QuotaRepository
}

func NewQuotaUsecase(
	db *gorm.DB,
	Quota_repository domain.QuotaRepository,
) *QuotaUsecase {
	return &QuotaUsecase{
		DB:              db,
		QuotaRepository: Quota_repository,
	}
}

func (a *QuotaUsecase) CreateQuota(ctx context.Context, request *model.WriteQuotaRequest) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	Quota := &domain.Quota{
		ScheduleID: request.ScheduleID,
		ClassID:    request.ClassID,
		Price:      request.Price,
		Quota:      request.Capacity,
		Capacity:   request.Capacity, // Assuming Capacity is the same as Quota
	}

	if err := a.QuotaRepository.Insert(tx, Quota); err != nil {
		return fmt.Errorf("failed to create Quota: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *QuotaUsecase) CreateQuotaBulk(ctx context.Context, requests []*model.WriteQuotaRequest) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	quotas := make([]*domain.Quota, len(requests))
	for i, request := range requests {
		quotas[i] = &domain.Quota{
			ScheduleID: request.ScheduleID,
			ClassID:    request.ClassID,
			Price:      request.Price,
			Quota:      request.Capacity,
			Capacity:   request.Capacity,
		}
	}

	if err := a.QuotaRepository.InsertBulk(tx, quotas); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create quotas in bulk: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *QuotaUsecase) ListQuotas(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadQuotaResponse, int, error) {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := a.QuotaRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count Quotas: %w", err)
	}

	Quotas, err := a.QuotaRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all Quotas: %w", err)
	}

	responses := make([]*model.ReadQuotaResponse, len(Quotas))
	for i, Quota := range Quotas {
		responses[i] = mapper.QuotaToResponse(Quota)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (a *QuotaUsecase) GetQuotaByID(ctx context.Context, id uint) (*model.ReadQuotaResponse, error) {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	Quota, err := a.QuotaRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get Quota: %w", err)
	}

	if Quota == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.QuotaToResponse(Quota), nil
}

func (a *QuotaUsecase) UpdateQuota(ctx context.Context, request *model.UpdateQuotaRequest) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing Quota
	Quota, err := a.QuotaRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find Quota: %w", err)
	}
	if Quota == nil {
		return errs.ErrNotFound
	}

	Quota.ScheduleID = request.ScheduleID
	Quota.ClassID = request.ClassID
	Quota.Capacity = request.Capacity // Assuming Capacity is the same as Quota
	Quota.Quota = request.Capacity
	Quota.Price = request.Price

	// Save updated Quota
	if err := a.QuotaRepository.Update(tx, Quota); err != nil {
		return fmt.Errorf("failed to update Quota: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *QuotaUsecase) DeleteQuota(ctx context.Context, id uint) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	Quota, err := a.QuotaRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get Quota: %w", err)
	}
	if Quota == nil {
		return errs.ErrNotFound
	}

	if err := a.QuotaRepository.Delete(tx, Quota); err != nil {
		return fmt.Errorf("failed to delete Quota: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
