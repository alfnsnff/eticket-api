package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"eticket-api/pkg/utils/tx"
	"fmt"

	"gorm.io/gorm"
)

type ShipUsecase struct {
	Tx             *tx.TxManager
	ShipRepository *repository.ShipRepository
}

func NewShipUsecase(
	tx *tx.TxManager,
	ship_repository *repository.ShipRepository,
) *ShipUsecase {
	return &ShipUsecase{
		Tx:             tx,
		ShipRepository: ship_repository,
	}
}

func (sh *ShipUsecase) CreateShip(ctx context.Context, request *model.WriteShipRequest) error {
	ship := mapper.ShipMapper.FromWrite(request)

	if ship.ShipName == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return sh.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return sh.ShipRepository.Create(tx, ship)
	})
}

func (sh *ShipUsecase) GetAllShips(ctx context.Context, limit, offset int) ([]*model.ReadShipResponse, int, error) {

	ships := []*entity.Ship{}
	var total int64

	err := sh.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = sh.ShipRepository.Count(tx)
		if err != nil {
			return err
		}
		ships, err = sh.ShipRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ShipMapper.ToModels(ships), int(total), nil
}

func (sh *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*model.ReadShipResponse, error) {
	ship := new(entity.Ship)

	err := sh.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		ship, err = sh.ShipRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}

	if ship == nil {
		return nil, errors.New("ship not found")
	}

	return mapper.ShipMapper.ToModel(ship), nil
}

func (sh *ShipUsecase) UpdateShip(ctx context.Context, id uint, request *model.UpdateShipRequest) error {
	ship := mapper.ShipMapper.FromUpdate(request)
	ship.ID = id

	if ship.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}

	if ship.ShipName == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return sh.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return sh.ShipRepository.Update(tx, ship)
	})

}

func (sh *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {

	return sh.Tx.Execute(ctx, func(tx *gorm.DB) error {
		ship, err := sh.ShipRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ship == nil {
			return errors.New("ship not found")
		}
		return sh.ShipRepository.Delete(tx, ship)
	})

}
