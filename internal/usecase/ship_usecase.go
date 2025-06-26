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

type ShipUsecase struct {
	DB             *gorm.DB
	ShipRepository domain.ShipRepository
}

func NewShipUsecase(
	db *gorm.DB,
	ship_repository domain.ShipRepository,
) *ShipUsecase {
	return &ShipUsecase{
		DB:             db,
		ShipRepository: ship_repository,
	}
}

func (sh *ShipUsecase) CreateShip(ctx context.Context, request *model.WriteShipRequest) error {
	tx := sh.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	ship := &domain.Ship{
		ShipName:      request.ShipName,
		ShipType:      request.ShipType,
		ShipAlias:     request.ShipAlias,
		Status:        request.Status,
		YearOperation: request.YearOperation,
		ImageLink:     request.ImageLink,
		Description:   request.Description,
	}

	if err := sh.ShipRepository.Insert(tx, ship); err != nil {
		return fmt.Errorf("failed to create ship: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sh *ShipUsecase) ListShips(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadShipResponse, int, error) {
	tx := sh.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	total, err := sh.ShipRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count ships: %w", err)
	}

	ships, err := sh.ShipRepository.FindAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all ships: %w", err)
	}

	responses := make([]*model.ReadShipResponse, len(ships))
	for i, ship := range ships {
		responses[i] = mapper.ShipToResponse(ship)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (sh *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*model.ReadShipResponse, error) {
	tx := sh.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	ship, err := sh.ShipRepository.FindByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ship: %w", err)
	}
	if ship == nil {
		return nil, errs.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.ShipToResponse(ship), nil
}

func (sh *ShipUsecase) UpdateShip(ctx context.Context, request *model.UpdateShipRequest) error {
	tx := sh.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Fetch existing allocation
	ship, err := sh.ShipRepository.FindByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find ship: %w", err)
	}
	if ship == nil {
		return errs.ErrNotFound
	}

	ship.ShipName = request.ShipName
	ship.ShipType = request.ShipType
	ship.ShipAlias = request.ShipAlias
	ship.Status = request.Status
	ship.YearOperation = request.YearOperation
	ship.ImageLink = request.ImageLink
	ship.Description = request.Description

	if err := sh.ShipRepository.Update(tx, ship); err != nil {
		return fmt.Errorf("failed to update ship: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (sh *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {
	tx := sh.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	ship, err := sh.ShipRepository.FindByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get ship: %w", err)
	}
	if ship == nil {
		return errs.ErrNotFound
	}

	if err := sh.ShipRepository.Delete(tx, ship); err != nil {
		return fmt.Errorf("failed to delete ship: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
