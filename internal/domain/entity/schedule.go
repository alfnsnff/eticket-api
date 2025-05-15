package entity

import "time"

type Schedule struct {
	ID                uint       `gorm:"column:id;primaryKey"`
	RouteID           uint       `gorm:"column:route_id;not null;index"`
	ShipID            uint       `gorm:"column:ship_id;not null;index"`
	DepartureDatetime *time.Time `gorm:"column:departure_datetime;"`
	ArrivalDatetime   *time.Time `gorm:"column:arrival_datetime;"`
	Status            *string    `gorm:"column:status;type:varchar(24);not null"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;not null"`

	Route Route `gorm:"foreignKey:RouteID" json:"route"` // Gorm will create the relationship
	Ship  Ship  `gorm:"foreignKey:ShipID" json:"ship"`   // Gorm will create the relationship
}

func (sch *Schedule) TableName() string {
	return "schedule"
}
