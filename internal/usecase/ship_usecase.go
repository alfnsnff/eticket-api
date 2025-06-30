package usecase

import (
	"context"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
)

type ShipUsecase struct {
	Transactor     transact.Transactor
	ShipRepository domain.ShipRepository
}

func NewShipUsecase(

	transactor transact.Transactor,
	ship_repository domain.ShipRepository,
) *ShipUsecase {
	return &ShipUsecase{

		Transactor:     transactor,
		ShipRepository: ship_repository,
	}
}

func (uc *ShipUsecase) CreateShip(ctx context.Context, e *domain.Ship) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ship := &domain.Ship{
			ShipName:      e.ShipName,
			ShipType:      e.ShipType,
			ShipAlias:     e.ShipAlias,
			Status:        e.Status,
			YearOperation: e.YearOperation,
			ImageLink:     e.ImageLink,
			Description:   e.Description,
		}
		if err := uc.ShipRepository.Insert(ctx, tx, ship); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to create ship: %w", err)
		}
		return nil
	})
}

func (uc *ShipUsecase) ListShips(ctx context.Context, limit, offset int, sort, search string) ([]*domain.Ship, int, error) {
	var err error
	var total int64
	var ships []*domain.Ship
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		total, err = uc.ShipRepository.Count(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to count ships: %w", err)
		}
		ships, err = uc.ShipRepository.FindAll(ctx, tx, limit, offset, sort, search)
		if err != nil {
			return fmt.Errorf("failed to get all ships: %w", err)
		}
		return nil
	}); err != nil {
		return nil, 0, fmt.Errorf("failed to list ships: %w", err)
	}
	return ships, int(total), nil
}

func (uc *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*domain.Ship, error) {
	var err error
	var ship *domain.Ship
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ship, err = uc.ShipRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get ship: %w", err)
		}
		if ship == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get ship by id: %w", err)
	}
	return ship, nil
}

func (uc *ShipUsecase) UpdateShip(ctx context.Context, e *domain.Ship) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ship, err := uc.ShipRepository.FindByID(ctx, tx, e.ID)
		if err != nil {
			return fmt.Errorf("failed to find ship: %w", err)
		}
		if ship == nil {
			return errs.ErrNotFound
		}

		ship.ShipName = e.ShipName
		ship.ShipType = e.ShipType
		ship.ShipAlias = e.ShipAlias
		ship.Status = e.Status
		ship.YearOperation = e.YearOperation
		ship.ImageLink = e.ImageLink
		ship.Description = e.Description

		if err := uc.ShipRepository.Update(ctx, tx, ship); err != nil {
			return fmt.Errorf("failed to update ship: %w", err)
		}

		return nil
	})
}

func (uc *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		ship, err := uc.ShipRepository.FindByID(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to get ship: %w", err)
		}
		if ship == nil {
			return errs.ErrNotFound
		}

		if err := uc.ShipRepository.Delete(ctx, tx, ship); err != nil {
			return fmt.Errorf("failed to delete ship: %w", err)
		}
		return nil
	})
}
