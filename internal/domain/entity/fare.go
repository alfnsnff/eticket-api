package entity

import "time"

type Fare struct {
	ID          uint      `gorm:"column:id;primaryKey"`
	RouteID     uint      `gorm:"column:route_id;not null;index;"`
	ManifestID  uint      `gorm:"column:manifest_id;not null;index;"`
	TicketPrice float32   `gorm:"column:ticket_price;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null"`

	Route    Route    `gorm:"foreignKey:RouteID"`
	Manifest Manifest `gorm:"foreignKey:ManifestID"`
}

func (f *Fare) TableName() string {
	return "fare"
}
