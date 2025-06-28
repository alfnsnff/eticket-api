package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
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
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Quota) error
	InsertBulk(ctx context.Context, conn gotann.Connection, Quotas []*Quota) error
	Update(ctx context.Context, conn gotann.Connection, entity *Quota) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, Quotas []*Quota) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Quota) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Quota, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Quota, error)
	FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*Quota, error)
	FindByScheduleIDAndClassID(ctx context.Context, conn gotann.Connection, scheduleID uint, classID uint) (*Quota, error)
}
