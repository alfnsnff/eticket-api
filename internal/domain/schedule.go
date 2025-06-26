package domain

import (
	"eticket-api/pkg/gotann"
	"time"

	"gorm.io/gorm"
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
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *Schedule) error
	Inserts(conn *gotann.Connection, entity *Schedule) error
	InsertBulk(db *gorm.DB, schedules []*Schedule) error
	Update(db *gorm.DB, entity *Schedule) error
	UpdateBulk(db *gorm.DB, schedules []*Schedule) error
	Delete(db *gorm.DB, entity *Schedule) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Schedule, error)
	FindByID(db *gorm.DB, id uint) (*Schedule, error)
	FindAllScheduled(db *gorm.DB) ([]*Schedule, error)
	FindActiveSchedule(db *gorm.DB) ([]*Schedule, error)
}
