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

func (f *FareUsecase) CreateFare(ctx context.Context, request *model.WriteFareRequest) error {
	fare := mapper.FareMapper.FromWrite(request)

	if fare.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if fare.ManifestID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if fare.Price == 0 {
		return fmt.Errorf("fare cannot be zero")
	}

	return tx.Execute(ctx, f.DB, func(tx *gorm.DB) error {
		return f.FareRepository.Create(tx, fare)
	})
}

func (f *FareUsecase) GetAllFares(ctx context.Context) ([]*model.ReadFareResponse, error) {
	fares := []*entity.Fare{}

	err := tx.Execute(ctx, f.DB, func(tx *gorm.DB) error {
		var err error
		fares, err = f.FareRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all Fares: %w", err)
	}

	return mapper.FareMapper.ToModels(fares), nil
}

func (f *FareUsecase) GetFareByID(ctx context.Context, id uint) (*model.ReadFareResponse, error) {
	fare := new(entity.Fare)

	err := tx.Execute(ctx, f.DB, func(tx *gorm.DB) error {
		var err error
		fare, err = f.FareRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get Fare by ID: %w", err)
	}

	if fare == nil {
		return nil, errors.New("fare not found")
	}

	return mapper.FareMapper.ToModel(fare), nil
}

func (f *FareUsecase) UpdateFare(ctx context.Context, id uint, request *model.UpdateFareRequest) error {
	fare := mapper.FareMapper.FromUpdate(request)
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
		return fmt.Errorf("fare cannot be zero")
	}

	return tx.Execute(ctx, f.DB, func(tx *gorm.DB) error {
		return f.FareRepository.Update(tx, fare)
	})
}

func (f *FareUsecase) DeleteFare(ctx context.Context, id uint) error {

	return tx.Execute(ctx, f.DB, func(tx *gorm.DB) error {
		fare, err := f.FareRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if fare == nil {
			return errors.New("fare not found")
		}
		return f.FareRepository.Delete(tx, fare)
	})

}
