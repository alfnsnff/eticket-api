package domain

import (
	"time"

	"gorm.io/gorm"
)

type Route struct {
	ID                uint      `gorm:"column:id;primaryKey" json:"id"`
	DepartureHarborID uint      `gorm:"column:departure_harbor_id;not null;index;"`
	ArrivalHarborID   uint      `gorm:"column:arrival_harbor_id;not null;index;"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`

	DepartureHarbor Harbor `gorm:"foreignKey:DepartureHarborID"` // Gorm will create the relationship
	ArrivalHarbor   Harbor `gorm:"foreignKey:ArrivalHarborID"`   // Gorm will create the relationship
}

func (r *Route) TableName() string {
	return "route"
}

type RouteRepository interface {
	Create(db *gorm.DB, entity *Route) error
	Update(db *gorm.DB, entity *Route) error
	Delete(db *gorm.DB, entity *Route) error
	Count(db *gorm.DB) (int64, error)
	GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*Route, error)
	GetByID(db *gorm.DB, id uint) (*Route, error)
}
