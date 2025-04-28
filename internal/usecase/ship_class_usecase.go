package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/repository"
	"fmt"
)

type ShipClassUsecase struct {
	ShipClassRepository *repository.ShipClassRepository
}

func NewShipClassUsecase(shipClassRepository *repository.ShipClassRepository) ShipClassUsecase {
	return ShipClassUsecase{ShipClassRepository: shipClassRepository}
}

func (s *ShipClassUsecase) CreateShipClass(shipClass *entities.ShipClass) error {
	if shipClass.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if shipClass.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}

	// Check if the combination already exists
	existing, err := s.ShipClassRepository.GetByShipAndClass(shipClass.ShipID, shipClass.ClassID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("ship class with this ship ID and class ID already exists")
	}

	return s.ShipClassRepository.Create(shipClass)
}

// GetAllshipes retrieves all ships
func (s *ShipClassUsecase) GetAllShipClasses() ([]*entities.ShipClass, error) {
	return s.ShipClassRepository.GetAll()
}

// GetshipByID retrieves a ship by its ID
func (s *ShipClassUsecase) GetShipClassByID(id uint) (*entities.ShipClass, error) {
	shipClass, err := s.ShipClassRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if shipClass == nil {
		return nil, errors.New("ship class not found")
	}
	return shipClass, nil
}

// GetshipByID retrieves a ship by its ID
func (s *ShipClassUsecase) GetShipClassByShipID(shipId uint) ([]*entities.ShipClass, error) {
	shipClasses, err := s.ShipClassRepository.GetByShipID(shipId)
	if err != nil {
		return nil, err
	}
	if shipClasses == nil {
		return nil, errors.New("ship class not found")
	}
	return shipClasses, nil
}

// Updateship updates an existing ship
func (s *ShipClassUsecase) UpdateShipClass(id uint, shipClass *entities.ShipClass) error {
	shipClass.ID = id

	if shipClass.ID == 0 {
		return fmt.Errorf("shipClass ID cannot be zero")
	}
	if shipClass.ShipID == 0 {
		return fmt.Errorf("shipClass ship ID cannot be zero")
	}
	if shipClass.ClassID == 0 {
		return fmt.Errorf("shipClass class ID cannot be zero")
	}
	return s.ShipClassRepository.Update(shipClass)
}

// Deleteship deletes a ship by its ID
func (s *ShipClassUsecase) DeleteShipClass(id uint) error {
	ship, err := s.ShipClassRepository.GetByID(id)
	if err != nil {
		return err
	}
	if ship == nil {
		return errors.New("ship class not found")
	}
	return s.ShipClassRepository.Delete(id)
}
