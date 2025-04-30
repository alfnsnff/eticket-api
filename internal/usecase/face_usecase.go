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

type FareUsecase struct {
	DB             *gorm.DB
	FareRepository *repository.FareRepository
}

func NewFareUsecase(db *gorm.DB, fare_repository *repository.FareRepository) *FareUsecase {
	return &FareUsecase{
		DB:             db,
		FareRepository: fare_repository,
	}
}

// CreateFare validates and creates a new Fare
func (s *FareUsecase) CreateFare(ctx context.Context, request *model.WriteFareRequest) error {
	fare := mapper.ToFareEntity(request)

	if fare.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if fare.ManifestID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if fare.Price == 0 {
		return fmt.Errorf("Fare cannot be zero")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.FareRepository.Create(tx, fare)
	})
}

// GetAllFares retrieves all Fares
func (s *FareUsecase) GetAllFares(ctx context.Context) ([]*model.ReadFareResponse, error) {
	fares := []*entities.Fare{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		fares, err = s.FareRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all Fares: %w", err)
	}

	return mapper.ToFaresModel(fares), nil
}

// GetFareByID retrieves a Fare by its ID
func (s *FareUsecase) GetFareByID(ctx context.Context, id uint) (*model.ReadFareResponse, error) {
	fare := new(entities.Fare)

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		fare, err = s.FareRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get Fare by ID: %w", err)
	}

	if fare == nil {
		return nil, errors.New("Fare not found")
	}

	return mapper.ToFareModel(fare), nil
}

// GetFareByRouteID retrieves Fares by route ID
func (s *FareUsecase) GetFareByRouteID(ctx context.Context, routeID uint) ([]*model.ReadFareResponse, error) {
	fares := []*entities.Fare{}

	err := tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		var err error
		fares, err = s.FareRepository.GetByRouteID(tx, routeID)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get Fare by route ID: %w", err)
	}

	if fares == nil {
		return nil, errors.New("Fares not found for this route")
	}

	return mapper.ToFaresModel(fares), nil
}

// UpdateFare updates an existing Fare
func (s *FareUsecase) UpdateFare(ctx context.Context, id uint, request *model.WriteFareRequest) error {
	fare := mapper.ToFareEntity(request)

	fare.ID = id

	if fare.ID == 0 {
		return fmt.Errorf("fare ID cannot be zero")
	}
	if fare.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if fare.ManifestID == 0 {
		return fmt.Errorf("manifest ID cannot be zero")
	}
	if fare.Price == 0 {
		return fmt.Errorf("Fare cannot be zero")
	}

	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		return s.FareRepository.Update(tx, fare)
	})
}

// DeleteFare deletes a Fare by its ID
func (s *FareUsecase) DeleteFare(ctx context.Context, id uint) error {
	return tx.Execute(ctx, s.DB, func(tx *gorm.DB) error {
		Fare, err := s.FareRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if Fare == nil {
			return errors.New("Fare not found")
		}
		return s.FareRepository.Delete(tx, Fare)
	})
}
