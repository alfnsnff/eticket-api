package allocation

import (
	"context"
	"errors"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
	"eticket-api/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type AllocationUsecase struct {
	Tx                   *tx.TxManager
	AllocationRepository *repository.AllocationRepository
	ScheduleRepository   *repository.ScheduleRepository
	FareRepository       *repository.FareRepository
}

func NewAllocationUsecase(
	tx *tx.TxManager,
	allocation_repository *repository.AllocationRepository,
	schedule_repository *repository.ScheduleRepository,
	fare_repository *repository.FareRepository,
) *AllocationUsecase {
	return &AllocationUsecase{
		Tx:                   tx,
		AllocationRepository: allocation_repository,
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
	return a.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return a.AllocationRepository.Create(tx, allocation)
	})
}

func (a *AllocationUsecase) GetAllAllocations(ctx context.Context, limit, offset int) ([]*model.ReadAllocationResponse, int, error) {
	allocations := []*entity.Allocation{}
	var total int64
	err := a.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		total, err = a.AllocationRepository.Count(tx)
		if err != nil {
			return err
		}
		allocations, err = a.AllocationRepository.GetAll(tx, limit, offset)
		return err
	})

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all allocations: %w", err)
	}

	return mapper.AllocationMapper.ToModels(allocations), int(total), nil
}

func (a *AllocationUsecase) GetAllocationByID(ctx context.Context, id uint) (*model.ReadAllocationResponse, error) {
	allocation := new(entity.Allocation)
	err := a.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		allocation, err = a.AllocationRepository.GetByID(tx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get allocation by ID: %w", err)
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

	return a.Tx.Execute(ctx, func(tx *gorm.DB) error {
		return a.AllocationRepository.Update(tx, allocation)
	})
}

func (a *AllocationUsecase) DeleteAllocation(ctx context.Context, id uint) error {

	return a.Tx.Execute(ctx, func(tx *gorm.DB) error {
		allocation, err := a.AllocationRepository.GetByID(tx, id)
		if err != nil {
			return err
		}
		if allocation == nil {
			return errors.New("allocation not found")
		}
		return a.AllocationRepository.Delete(tx, allocation)
	})

}
