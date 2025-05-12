package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ShipUsecase struct {
	DB             *gorm.DB
	ShipRepository *repository.ShipRepository
}

func NewShipUsecase(db *gorm.DB, ship_repository *repository.ShipRepository) *ShipUsecase {
	return &ShipUsecase{DB: db, ShipRepository: ship_repository}
}

func (sh *ShipUsecase) CreateShip(ctx context.Context, request *model.WriteShipRequest) error {
	ship := mapper.ShipMapper.FromWrite(request)

	if ship.ShipName == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return tx.Execute(ctx, sh.DB, func(tx *gorm.DB) error {
		return sh.ShipRepository.Create(tx, ship)
	})
}

func (sh *ShipUsecase) GetAllShips(ctx context.Context) ([]*model.ReadShipResponse, error) {

	ships := []*entity.Ship{}

	err := tx.Execute(ctx, sh.DB, func(tx *gorm.DB) error {
		var err error
		ships, err = sh.ShipRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ShipMapper.ToModels(ships), nil
}

func (sh *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*model.ReadShipResponse, error) {
	ship := new(entity.Ship)

	err := tx.Execute(ctx, sh.DB, func(tx *gorm.DB) error {
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

	return tx.Execute(ctx, sh.DB, func(tx *gorm.DB) error {
		return sh.ShipRepository.Update(tx, ship)
	})

}

func (sh *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {

	return tx.Execute(ctx, sh.DB, func(tx *gorm.DB) error {
		ship, err := sh.ShipRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ship == nil {
			return errors.New("route not found")
		}
		return sh.ShipRepository.Delete(tx, ship)
	})

}
