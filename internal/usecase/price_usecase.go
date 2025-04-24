package usecase

import (
	"errors"
	"eticket-api/internal/domain/entities"
	"fmt"
)

type PriceUsecase struct {
	PriceRepository entities.PriceRepositoryInterface
}

func NewPriceUsecase(priceRepository entities.PriceRepositoryInterface) PriceUsecase {
	return PriceUsecase{PriceRepository: priceRepository}
}

// Createship validates and creates a new ship
func (s *PriceUsecase) CreatePrice(price *entities.Price) error {
	if price.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if price.ShipClassID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if price.Price == 0 {
		return fmt.Errorf("price cannot be zero")
	}
	return s.PriceRepository.Create(price)
}

// GetAllshipes retrieves all ships
func (s *PriceUsecase) GetAllPrices() ([]*entities.Price, error) {
	return s.PriceRepository.GetAll()
}

// GetshipByID retrieves a ship by its ID
func (s *PriceUsecase) GetPriceByID(id uint) (*entities.Price, error) {
	ship, err := s.PriceRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if ship == nil {
		return nil, errors.New("price not found")
	}
	return ship, nil
}

// GetshipByID retrieves a ship by its ID
func (s *PriceUsecase) GetPriceByRouteID(id uint) ([]*entities.Price, error) {
	ship, err := s.PriceRepository.GetByRouteID(id)
	if err != nil {
		return nil, err
	}
	if ship == nil {
		return nil, errors.New("price not found")
	}
	return ship, nil
}

// Updateship updates an existing ship
func (s *PriceUsecase) UpdatePrice(id uint, price *entities.Price) error {
	price.ID = id

	if price.ID == 0 {
		return fmt.Errorf("price ID cannot be zero")
	}
	if price.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if price.ShipClassID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if price.Price == 0 {
		return fmt.Errorf("price cannot be zero")
	}
	return s.PriceRepository.Update(price)
}

// Deleteship deletes a ship by its ID
func (s *PriceUsecase) DeletePrice(id uint) error {
	ship, err := s.PriceRepository.GetByID(id)
	if err != nil {
		return err
	}
	if ship == nil {
		return errors.New("price not found")
	}
	return s.PriceRepository.Delete(id)
}
