package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"fmt"
)

type ShipUsecase struct {
	ShipRepository entities.ShipRepositoryInterface
}

func NewShipUsecase(shipRepository entities.ShipRepositoryInterface) ShipUsecase {
	return ShipUsecase{ShipRepository: shipRepository}
}

// Createship validates and creates a new ship
func (s *ShipUsecase) CreateShip(ship *entities.Ship) error {
	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}
	return s.ShipRepository.Create(ship)
}

// GetAllshipes retrieves all ships
func (s *ShipUsecase) GetAllShips() ([]*entities.Ship, error) {
	return s.ShipRepository.GetAll()
}

// GetshipByID retrieves a ship by its ID
func (s *ShipUsecase) GetShipByID(id uint) (*entities.Ship, error) {
	ship, err := s.ShipRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ship == nil {
		return nil, errors.New("ship not found")
	}
	return ship, nil
}

// Updateship updates an existing ship
func (s *ShipUsecase) UpdateShip(id uint, ship *entities.Ship) error {
	ship.ID = id

	if ship.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}
	if ship.Name == "" {
		return fmt.Errorf("ship name cannot be empty")
	}
	return s.ShipRepository.Update(ship)
}

// Deleteship deletes a ship by its ID
func (s *ShipUsecase) DeleteShip(id uint) error {
	ship, err := s.ShipRepository.GetByID(id)
	if err != nil {
		return err
	}
	if ship == nil {
		return errors.New("ship not found")
	}
	return s.ShipRepository.Delete(id)
}
