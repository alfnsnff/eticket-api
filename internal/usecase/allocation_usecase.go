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

type AllocationUsecase struct {
	DB                   *gorm.DB
	allocationRepository *repository.AllocationRepository
	ScheduleRepository   *repository.ScheduleRepository
	FareRepository       *repository.FareRepository
}

func NewAllocationUsecase(
	db *gorm.DB,
	allocation_repository *repository.AllocationRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
) *AllocationUsecase {
	return &AllocationUsecase{
		DB:                   db,
		allocationRepository: allocation_repository,
		ScheduleRepository:   schedule_repository,
		FareRepository:       fare_repository,
	}
}

func (a *AllocationUsecase) CreateAllocation(ctx context.Context, request *model.WriteAllocationRequest) error {
	allocation := mapper.AllocationMapper.FromWrite(request)

	if allocation.ScheduleID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if allocation.ClassID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}
	return tx.Execute(ctx, a.DB, func(tx *gorm.DB) error {
		return a.allocationRepository.Create(tx, allocation)
	})
}

func (a *AllocationUsecase) GetAllAllocations(ctx context.Context) ([]*model.ReadAllocationResponse, error) {
	allocations := []*entity.Allocation{}

	err := tx.Execute(ctx, a.DB, func(tx *gorm.DB) error {
		var err error
		allocations, err = a.allocationRepository.GetAll(tx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get all allocations: %w", err)
	}

	return mapper.AllocationMapper.ToModels(allocations), nil
}

func (a *AllocationUsecase) GetAllocationByID(ctx context.Context, id uint) (*model.ReadAllocationResponse, error) {
	allocation := new(entity.Allocation)
	err := tx.Execute(ctx, a.DB, func(tx *gorm.DB) error {
		var err error
		allocation, err = a.allocationRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get Allocation by ID: %w", err)
	}

	if allocation == nil {
		return nil, errors.New("allocation not found")
	}

	return mapper.AllocationMapper.ToModel(allocation), nil
}

func (a *AllocationUsecase) UpdateAllocation(ctx context.Context, id uint, request *model.UpdateAllocationRequest) error {
	allocation := mapper.AllocationMapper.FromUpdate(request)
	allocation.ID = id

	if allocation.ID == 0 {
		return fmt.Errorf("allocation ID cannot be zero")
	}

	if allocation.ScheduleID == 0 {
		return fmt.Errorf("schedule ID cannot be zero")
	}

	if allocation.ClassID == 0 {
		return fmt.Errorf("class ID cannot be zero")
	}

	if allocation.Quota == 0 {
		return fmt.Errorf("quota name cannot be zero")
	}

	return tx.Execute(ctx, a.DB, func(tx *gorm.DB) error {
		return a.allocationRepository.Update(tx, allocation)
	})
}

func (a *AllocationUsecase) DeleteAllocation(ctx context.Context, id uint) error {

	return tx.Execute(ctx, a.DB, func(tx *gorm.DB) error {
		allocation, err := a.allocationRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if allocation == nil {
			return errors.New("allocation not found")
		}
		return a.allocationRepository.Delete(tx, allocation)
	})

}
