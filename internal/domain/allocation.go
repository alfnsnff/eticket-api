package domain

import (
	"time"

	"gorm.io/gorm"
)

type Allocation struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ScheduleID uint      `gorm:"column:schedule_id;not null"`
	ClassID    uint      `gorm:"column:class_id;not null"`
	Quota      int       `gorm:"column:quota;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"`
}

func (a *Allocation) TableName() string {
	return "allocation"
}

type AllocationRepository interface {
	Create(db *gorm.DB, entity *Allocation) error
	Update(db *gorm.DB, entity *Allocation) error
	Delete(db *gorm.DB, entity *Allocation) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Allocation, error)
	GetByID(db *gorm.DB, id uint) (*Allocation, error)
	LockByScheduleAndClass(db *gorm.DB, scheduleID uint, classID uint) (*Allocation, error)
	GetByScheduleAndClass(db *gorm.DB, scheduleID uint, classID uint) (*Allocation, error)
	FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*Allocation, error)
	CreateBulk(db *gorm.DB, allocations []*Allocation) error
}

type AllocationUsecase interface {
}
