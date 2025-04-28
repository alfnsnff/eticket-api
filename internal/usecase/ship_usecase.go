package usecase

import (
	"context"
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"
	"fmt"

	"gorm.io/gorm"
)

type ShipUsecase struct {
	DB             *gorm.DB
	ShipRepository *repository.ShipRepository
}

func NewShipUsecase(db *gorm.DB, shipRepository *repository.ShipRepository) *ShipUsecase {
	return &ShipUsecase{DB: db, ShipRepository: shipRepository}
}

// Createship validates and creates a new ship
func (s *ShipUsecase) CreateShip(ctx context.Context, ship *entities.Ship) error {
	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ShipRepository.Create(txDB, ship)
	})
}

// GetAllshipes retrieves all ships
func (s *ShipUsecase) GetAllShips(ctx context.Context) ([]*entities.Ship, error) {

	var ships []*entities.Ship

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		ships, err = s.ShipRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all books: %w", err)
	}

	return ships, nil
}

// GetshipByID retrieves a ship by its ID
func (s *ShipUsecase) GetShipByID(ctx context.Context, id uint) (*entities.Ship, error) {

	var ship *entities.Ship

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		ship, err = s.ShipRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, err
	}
	if ship == nil {
		return nil, errors.New("ship not found")
	}

	return ship, nil
}

// Updateship updates an existing ship
func (s *ShipUsecase) UpdateShip(ctx context.Context, id uint, ship *entities.Ship) error {
	ship.ID = id

	if ship.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}

	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.ShipRepository.Update(txDB, ship)
	})

}

// Deleteship deletes a ship by its ID
func (s *ShipUsecase) DeleteShip(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		ship, err := s.ShipRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if ship == nil {
			return errors.New("route not found")
		}
		return s.ShipRepository.Delete(txDB, id)
	})
}
