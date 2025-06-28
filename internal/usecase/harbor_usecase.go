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

type HarborUsecase struct {
	Transactor       *transact.Transactor
	HarborRepository domain.HarborRepository
}

func NewHarborUsecase(

	transactor *transact.Transactor,
	harborRepository domain.HarborRepository,
) *HarborUsecase {
	return &HarborUsecase{

		Transactor:       transactor,
		HarborRepository: harborRepository,
	}
}

func (uc *HarborUsecase) CreateHarbor(ctx context.Context, request *model.WriteHarborRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		harbor := &domain.Harbor{
			HarborName:    request.HarborName,
			Status:        request.Status,
			HarborAlias:   request.HarborAlias,
			YearOperation: request.YearOperation,
		}
		if err := uc.HarborRepository.Insert(ctx, tx, harbor); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create harbor: %w", err)
		}
		return nil
	})
}

func (uc *HarborUsecase) ListHarbors(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadHarborResponse, int, error) {
	var err error
	var total int64
	var harbors []*domain.Harbor
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.HarborRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count harbors: %w", err)
		}
		harbors, err = uc.HarborRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get harbors: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list harbors: %w", err)
	}
	responses := make([]*model.ReadHarborResponse, len(harbors))
	for i, harbor := range harbors {
		responses[i] = mapper.HarborToResponse(harbor)
	}

	return responses, int(total), nil
}

func (uc *HarborUsecase) GetHarborByID(ctx context.Context, id uint) (*model.ReadHarborResponse, error) {
	var err error
	var harbor *domain.Harbor
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		harbor, err := uc.HarborRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get harbor: %w", err)
		}
		if harbor == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get harbor by ID: %w", err)
	}
	return mapper.HarborToResponse(harbor), nil
}

func (uc *HarborUsecase) UpdateHarbor(ctx context.Context, request *model.UpdateHarborRequest) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		harbor, err := uc.HarborRepository.FindByID(ctx, tx, request.ID)
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

		if err := uc.HarborRepository.Update(ctx, tx, harbor); err != nil {
			return fmt.Errorf("failed to update harbor: %w", err)
		}
		return nil
	})
}

func (uc *HarborUsecase) DeleteHarbor(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		harbor, err := uc.HarborRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get harbor: %w", err)
		}
		if harbor == nil {
			return errs.ErrNotFound
		}

		if err := uc.HarborRepository.Delete(ctx, tx, harbor); err != nil {
			return fmt.Errorf("failed to delete harbor: %w", err)
		}
		return nil
	})
}
