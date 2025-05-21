package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AllocationRepository struct {
	Repository[entity.Allocation]
}

func NewAllocationRepository() *AllocationRepository {
	return &AllocationRepository{}
}

func (ar *AllocationRepository) Count(db *gorm.DB) (int64, error) {
	allocations := []*entity.Allocation{}
	var total int64
	result := db.Find(&allocations).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (ar *AllocationRepository) GetAll(db *gorm.DB, limit, offset int) ([]*entity.Allocation, error) {
	allocations := []*entity.Allocation{}
	result := db.Limit(limit).Offset(offset).Find(&allocations)
	if result.Error != nil {
		return nil, result.Error
	}
	return allocations, nil
}

func (ar *AllocationRepository) GetByID(db *gorm.DB, id uint) (*entity.Allocation, error) {
	allocation := new(entity.Allocation)
	result := db.First(&allocation, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return allocation, result.Error
}

func (ar *AllocationRepository) LockByScheduleAndClass(db *gorm.DB, scheduleID uint, classID uint) (*entity.Allocation, error) {
	allocation := new(entity.Allocation)
	result := db.Where("schedule_id = ? AND class_id = ?", scheduleID, classID).Clauses(clause.Locking{Strength: "UPDATE"}).First(allocation)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return allocation, result.Error
}

func (ar *AllocationRepository) GetBySchedlueAndClass(db *gorm.DB, scheduleID uint, classID uint) (*entity.Allocation, error) {
	allocation := new(entity.Allocation)
	result := db.Where("schedule_id = ? AND class_id = ?", scheduleID, classID).Find(allocation)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return allocation, result.Error
}

func (ar *AllocationRepository) FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*entity.Allocation, error) {
	allocations := []*entity.Allocation{} // Declare an empty slice of pointers

	// CORRECT LINE: Pass a POINTER to the slice (&allocations)
	result := db.Where("schedule_id = ?", scheduleID).Find(&allocations)

	// Check for any errors that are NOT ErrRecordNotFound (Find doesn't return it)
	if result.Error != nil {
		return nil, result.Error // Handle database errors
	}

	// Return the slice (it will be empty if no records were found, and result.Error will be nil)
	return allocations, nil
}

func (tr *AllocationRepository) CreateBulk(db *gorm.DB, allocations []*entity.Allocation) error {
	result := db.Create(&allocations)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
