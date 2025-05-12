package entity

import "time"

type Schedule struct {
	ID               uint      `gorm:"column:id;primaryKey"`
	RouteID          uint      `gorm:"column:route_id;not null;index"`
	ShipID           uint      `gorm:"column:ship_id;not null;index"`
	ScheduleDatetime time.Time `gorm:"column:schedule_datetime;not null"`
	Status           *string   `gorm:"column:status;type:varchar(24);not null"`
	CreatedAt        time.Time `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null"`
}

func (sch *Schedule) TableName() string {
	return "schedule"
}
