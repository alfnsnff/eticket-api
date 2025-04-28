package entities

import "time"

type Ship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ShipClasses []ShipClass `gorm:"foreignKey:ShipID" json:"Classes"`
}
