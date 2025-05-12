package entity

import "time"

type Route struct {
	ID                uint      `gorm:"column:id;primaryKey" json:"id"`
	DepartureHarborID uint      `gorm:"column:departure_harbor_id;not null;index;" json:"departure_harbor_id"`
	ArrivalHarborID   uint      `gorm:"column:arrival_harbor_id;not null;index;" json:"arrival_harbor_id"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

func (r *Route) TableName() string {
	return "route"
}
