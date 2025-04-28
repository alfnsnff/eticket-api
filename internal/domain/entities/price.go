package entities

import "time"

type Price struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RouteID     uint      `json:"route_id"`
	ShipClassID uint      `json:"ship_class_id"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ShipClass ShipClass `gorm:"foreignKey:ShipClassID" json:"ship_class"`
}
