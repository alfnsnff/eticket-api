package fare

import (
	"context"
	"errors"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type FareUsecase struct {
	DB             *gorm.DB
	FareRepository *repository.FareRepository
}

func NewFareUsecase(
	db *gorm.DB,
	fareRepository *repository.FareRepository,
) *FareUsecase {
	return &FareUsecase{
		DB:             db,
		FareRepository: fareRepository,
	}
}

func (f *FareUsecase) CreateFare(ctx context.Context, request *model.WriteFareRequest) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	fare := &entity.Fare{
		RouteID:     request.RouteID,
		ManifestID:  request.ManifestID,
		TicketPrice: request.TicketPrice,
	}

	if fare.RouteID == 0 {
		return fmt.Errorf("route ID cannot be zero")
	}
	if fare.ManifestID == 0 {
		return fmt.Errorf("ship class ID cannot be zero")
	}
	if fare.TicketPrice == 0 {
		return fmt.Errorf("fare cannot be zero")
	}

	if err := f.FareRepository.Create(tx, fare); err != nil {
		return fmt.Errorf("failed to create fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (f *FareUsecase) GetAllFares(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadFareResponse, int, error) {
	tx := f.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := f.FareRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count fares: %w", err)
	}

	fares, err := f.FareRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get fares: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.FareMapper.ToModels(fares), int(total), nil
}

func (f *FareUsecase) GetFareByID(ctx context.Context, id uint) (*model.ReadFareResponse, error) {
	tx := f.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	fare, err := f.FareRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fare: %w", err)
	}
	if fare == nil {
		return nil, errors.New("fare not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.FareMapper.ToModel(fare), nil
}

func (f *FareUsecase) UpdateFare(ctx context.Context, request *model.UpdateFareRequest) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	if request.ID == 0 {
		return fmt.Errorf("fare ID cannot be zero")
	}

	// Fetch existing allocation
	fare, err := f.FareRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find fare: %w", err)
	}
	if fare == nil {
		return errors.New("fare not found")
	}

	fare.ManifestID = request.ManifestID
	fare.RouteID = request.RouteID
	fare.TicketPrice = request.TicketPrice

	if err := f.FareRepository.Update(tx, fare); err != nil {
		return fmt.Errorf("failed to update fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (f *FareUsecase) DeleteFare(ctx context.Context, id uint) error {
	tx := f.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	fare, err := f.FareRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get fare: %w", err)
	}
	if fare == nil {
		return errors.New("fare not found")
	}

	if err := f.FareRepository.Delete(tx, fare); err != nil {
		return fmt.Errorf("failed to delete fare: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
