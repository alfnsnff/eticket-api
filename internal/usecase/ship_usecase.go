package usecase

import (
	"errors"
	"eticket-api/internal/domain"
	"fmt"
)

type ShipUsecase struct {
	ShipRepository domain.ShipRepositoryInterface
}

func NewShipUsecase(shipRepository domain.ShipRepositoryInterface) ShipUsecase {
	return ShipUsecase{ShipRepository: shipRepository}
}

// Createship validates and creates a new ship
func (s *ShipUsecase) CreateShip(ship *domain.Ship) error {
	if ship.ShipName == "" {
		return fmt.Errorf("ship name cannot be empty")
	}
	return s.ShipRepository.Create(ship)
}

// GetAllshipes retrieves all ships
func (s *ShipUsecase) GetAllShips() ([]*domain.Ship, error) {
	return s.ShipRepository.GetAll()
}

// GetshipByID retrieves a ship by its ID
func (s *ShipUsecase) GetShipByID(id uint) (*domain.Ship, error) {
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
func (s *ShipUsecase) UpdateShip(ship *domain.Ship) error {
	if ship.ID == 0 {
		return fmt.Errorf("ship ID cannot be zero")
	}
	if ship.ShipName == "" {
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
