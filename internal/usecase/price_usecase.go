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

type PriceUsecase struct {
	DB              *gorm.DB
	PriceRepository *repository.PriceRepository
}

func NewPriceUsecase(db *gorm.DB, priceRepository *repository.PriceRepository) *PriceUsecase {
	return &PriceUsecase{
		DB:              db,
		PriceRepository: priceRepository,
	}
}

// CreatePrice validates and creates a new price
func (s *PriceUsecase) CreatePrice(ctx context.Context, price *entities.Price) error {
	if price.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if price.ShipClassID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if price.Price == 0 {
		return fmt.Errorf("price cannot be zero")
	}

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.PriceRepository.Create(txDB, price)
	})
}

// GetAllPrices retrieves all prices
func (s *PriceUsecase) GetAllPrices(ctx context.Context) ([]*entities.Price, error) {
	var prices []*entities.Price

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		prices, err = s.PriceRepository.GetAll(txDB)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all prices: %w", err)
	}

	return prices, nil
}

// GetPriceByID retrieves a price by its ID
func (s *PriceUsecase) GetPriceByID(ctx context.Context, id uint) (*entities.Price, error) {
	var price *entities.Price

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		price, err = s.PriceRepository.GetByID(txDB, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get price by ID: %w", err)
	}

	if price == nil {
		return nil, errors.New("price not found")
	}

	return price, nil
}

// GetPriceByRouteID retrieves prices by route ID
func (s *PriceUsecase) GetPriceByRouteID(ctx context.Context, routeID uint) ([]*entities.Price, error) {
	var prices []*entities.Price

	err := tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		var err error
		prices, err = s.PriceRepository.GetByRouteID(txDB, routeID)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get price by route ID: %w", err)
	}

	if prices == nil {
		return nil, errors.New("prices not found for this route")
	}

	return prices, nil
}

// UpdatePrice updates an existing price
func (s *PriceUsecase) UpdatePrice(ctx context.Context, id uint, price *entities.Price) error {
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

	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		return s.PriceRepository.Update(txDB, price)
	})
}

// DeletePrice deletes a price by its ID
func (s *PriceUsecase) DeletePrice(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(txDB *gorm.DB) error {
		price, err := s.PriceRepository.GetByID(txDB, id)
		if err != nil {
			return err
		}
		if price == nil {
			return errors.New("price not found")
		}
		return s.PriceRepository.Delete(txDB, id)
	})
}
