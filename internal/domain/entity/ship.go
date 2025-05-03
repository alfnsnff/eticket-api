package entity

import "time"

type Ship struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Status      string    `json:"status"`
	ShipType    string    `json:"ship_type"`
	Year        string    `json:"year"`
	Image       string    `json:"image"`
	Description string    `json:"Description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
