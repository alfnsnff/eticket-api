package domain

import (
	"time"

	"gorm.io/gorm"
)

type Quota struct {
	ID         uint      `gorm:"column:id;primaryKey"`
	ScheduleID uint      `gorm:"column:schedule_id;not null;uniqueIndex:idx_schedule_class"`
	ClassID    uint      `gorm:"column:class_id;not null;uniqueIndex:idx_schedule_class"`
	Quota      int       `gorm:"column:quota;not null"`
	Capacity   int       `gorm:"column:capacity;not null"`
	Price      float64   `gorm:"column:price;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`

	Class    Class    `gorm:"foreignKey:ClassID"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID"`
}

func (q *Quota) TableName() string {
	return "quota"
}

type QuotaRepository interface {
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *Quota) error
	InsertBulk(db *gorm.DB, Quotas []*Quota) error
	Update(db *gorm.DB, entity *Quota) error
	UpdateBulk(db *gorm.DB, Quotas []*Quota) error
	Delete(db *gorm.DB, entity *Quota) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Quota, error)
	FindByID(db *gorm.DB, id uint) (*Quota, error)
	FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*Quota, error)
	FindByScheduleIDAndClassID(db *gorm.DB, scheduleID uint, classID uint) (*Quota, error)
}
