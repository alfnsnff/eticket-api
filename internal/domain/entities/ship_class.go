package entities

import "time"

type ShipClass struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShipID    uint      `gorm:"not null" json:"ship_id"`
	ClassID   uint      `gorm:"not null" json:"class_id"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Ship  Ship  `gorm:"foreignKey:ShipID" json:"ship"`
	Class Class `gorm:"foreignKey:ClassID" json:"class"`
}
