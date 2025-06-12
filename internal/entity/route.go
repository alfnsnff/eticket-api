package entity

import (
	"time"
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
