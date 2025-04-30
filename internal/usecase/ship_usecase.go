package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
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

// Createship validates and creates a new ship
func (s *ShipUsecase) CreateShip(ctx context.Context, request *model.WriteShipRequest) error {
	ship := mapper.ToShipEntity(request)

	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ShipRepository.Create(tx, ship)
	})
}

// GetAllshipes retrieves all ships
func (s *ShipUsecase) GetAllShips(ctx context.Context) ([]*model.ReadShipResponse, error) {

	ships := []*entities.Ship{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		ships, err = s.ShipRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return mapper.ToShipsModel(ships), nil
}

// GetshipByID retrieves a ship by its ID
func (s *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*model.ReadShipResponse, error) {

	ship := new(entities.Ship)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		ship, err = s.ShipRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if ship == nil {
		return nil, errors.New("ship not found")
	}

	return mapper.ToShipModel(ship), nil
}

// Updateship updates an existing ship
func (s *ShipUsecase) UpdateShip(ctx context.Context, id uint, request *model.WriteShipRequest) error {
	ship := mapper.ToShipEntity(request)
	ship.ID = id

	if ship.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}

	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.ShipRepository.Update(tx, ship)
	})

}

// Deleteship deletes a ship by its ID
func (s *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		ship, err := s.ShipRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if ship == nil {
			return errors.New("route not found")
		}
		return s.ShipRepository.Delete(tx, ship)
	})
}
