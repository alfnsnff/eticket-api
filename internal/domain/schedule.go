package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type Schedule struct {
	ID                uint      `gorm:"column:id;primaryKey"`
	ShipID            uint      `gorm:"column:ship_id;not null;index"`
	DepartureHarborID uint      `gorm:"column:departure_harbor_id;not null;index;"`
	ArrivalHarborID   uint      `gorm:"column:arrival_harbor_id;not null;index;"`
	DepartureDatetime time.Time `gorm:"column:departure_datetime;;not null"`
	ArrivalDatetime   time.Time `gorm:"column:arrival_datetime;;not null"`
	Status            string    `gorm:"column:status;type:varchar(24);not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`

	Ship            Ship           `gorm:"foreignKey:ShipID"` // Gorm will create the relationship
	DepartureHarbor Harbor         `gorm:"foreignKey:DepartureHarborID"`
	ArrivalHarbor   Harbor         `gorm:"foreignKey:ArrivalHarborID"`
	Quotas          []Quota        `gorm:"foreignKey:ScheduleID"`
	ClaimSessions   []ClaimSession `gorm:"foreignKey:ScheduleID"`
}

func (sch *Schedule) TableName() string {
	return "schedule"
}

type ScheduleRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *Schedule) error
	InsertBulk(ctx context.Context, conn gotann.Connection, schedules []*Schedule) error
	Update(ctx context.Context, conn gotann.Connection, entity *Schedule) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, schedules []*Schedule) error
	Delete(ctx context.Context, conn gotann.Connection, entity *Schedule) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*Schedule, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*Schedule, error)
	FindActiveSchedules(ctx context.Context, conn gotann.Connection) ([]*Schedule, error)
}
