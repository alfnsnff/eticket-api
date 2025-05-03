package entity

import "time"

type Manifest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShipID    uint      `gorm:"not null;index;" json:"ship_id"`
	ClassID   uint      `gorm:"not null;index;" json:"class_id"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
