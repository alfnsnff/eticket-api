package domain

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID                uint       `gorm:"column:id;primaryKey"`
	RouteID           uint       `gorm:"column:route_id;not null;index"`
	ShipID            uint       `gorm:"column:ship_id;not null;index"`
	DepartureDatetime *time.Time `gorm:"column:departure_datetime;"`
	ArrivalDatetime   *time.Time `gorm:"column:arrival_datetime;"`
	Status            *string    `gorm:"column:status;type:varchar(24);not null"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null"`

	Route         Route          `gorm:"foreignKey:RouteID"` // Gorm will create the relationship
	Ship          Ship           `gorm:"foreignKey:ShipID"`  // Gorm will create the relationship
	ClaimSessions []ClaimSession `gorm:"foreignKey:ScheduleID"`
}

func (sch *Schedule) TableName() string {
	return "schedule"
}

type ScheduleRepository interface {
	Create(db *gorm.DB, entity *Schedule) error
	Update(db *gorm.DB, entity *Schedule) error
	Delete(db *gorm.DB, entity *Schedule) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Schedule, error)
	GetByID(db *gorm.DB, id uint) (*Schedule, error)
	GetAllScheduled(db *gorm.DB) ([]*Schedule, error)
	GetActiveSchedule(db *gorm.DB) ([]*Schedule, error)
}
