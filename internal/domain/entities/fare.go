package entities

import "time"

type Fare struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RouteID    uint      `json:"route_id"`
	ManifestID uint      `json:"manifest_id"`
	Price      float32   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Manifest Manifest `gorm:"foreignKey:ManifestID" json:"ship_class"`
}
