package allocation

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type AllocationUsecase struct {
	DB                   *gorm.DB // Assuming you have a DB field for the transaction manager
	AllocationRepository AllocationRepository
	FareRepository       FareRepository
}

func NewAllocationUsecase(
	db *gorm.DB,
	allocation_repository AllocationRepository,
	fare_repository FareRepository,
) *AllocationUsecase {
	return &AllocationUsecase{
		DB:                   db,
		AllocationRepository: allocation_repository,
		FareRepository:       fare_repository,
	}
}

func (a *AllocationUsecase) CreateAllocation(ctx context.Context, request *model.WriteAllocationRequest) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	allocation := &domain.Allocation{
		ScheduleID: request.ScheduleID,
		ClassID:    request.ClassID,
		Quota:      request.Quota,
	}

	if err := a.AllocationRepository.Create(tx, allocation); err != nil {
		return fmt.Errorf("failed to create allocation: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *AllocationUsecase) GetAllAllocations(ctx context.Context, limit, offset int, sort, search string) ([]*model.ReadAllocationResponse, int, error) {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	total, err := a.AllocationRepository.Count(tx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count allocations: %w", err)
	}

	allocations, err := a.AllocationRepository.GetAll(tx, limit, offset, sort, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all allocations: %w", err)
	}

	responses := make([]*model.ReadAllocationResponse, len(allocations))
	for i, allocation := range allocations {
		responses[i] = AllocationToResponse(allocation)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return responses, int(total), nil
}

func (a *AllocationUsecase) GetAllocationByID(ctx context.Context, id uint) (*model.ReadAllocationResponse, error) {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	allocation, err := a.AllocationRepository.GetByID(tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocation: %w", err)
	}

	if allocation == nil {
		return nil, errors.New("allocation not found")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return AllocationToResponse(allocation), nil
}

func (a *AllocationUsecase) UpdateAllocation(ctx context.Context, request *model.UpdateAllocationRequest) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	// Fetch existing allocation
	allocation, err := a.AllocationRepository.GetByID(tx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to find allocation: %w", err)
	}
	if allocation == nil {
		return errors.New("allocation not found")
	}

	allocation.ScheduleID = request.ScheduleID
	allocation.ClassID = request.ClassID
	allocation.Quota = request.Quota

	// Save updated allocation
	if err := a.AllocationRepository.Update(tx, allocation); err != nil {
		return fmt.Errorf("failed to update allocation: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (a *AllocationUsecase) DeleteAllocation(ctx context.Context, id uint) error {
	tx := a.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else {
			tx.Rollback()
		}
	}()

	allocation, err := a.AllocationRepository.GetByID(tx, id)
	if err != nil {
		return fmt.Errorf("failed to get allocation: %w", err)
	}
	if allocation == nil {
		return fmt.Errorf("allocation not found")
	}

	if err := a.AllocationRepository.Delete(tx, allocation); err != nil {
		return fmt.Errorf("failed to delete allocation: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
