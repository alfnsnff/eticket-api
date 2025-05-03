package entity

import "time"

type Fare struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RouteID    uint      `gorm:"not null;index;" json:"route_id"`
	ManifestID uint      `gorm:"not null;index;" json:"manifest_id"`
	Price      float32   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
