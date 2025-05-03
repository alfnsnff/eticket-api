package entity

import "time"

type Route struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	DepartureHarborID uint      `gorm:"not null;index;" json:"departure_harbor_id"`
	ArrivalHarborID   uint      `gorm:"not null;index;" json:"arrival_harbor_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
