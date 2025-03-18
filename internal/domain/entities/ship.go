package entities

import "time"

type Ship struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ShipName  string    `gorm:"not null" json:"ship_name"`
	Capacity  uint      `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ShipRepositoryInterface interface {
	Create(ship *Ship) error
	GetAll() ([]*Ship, error)
	GetByID(id uint) (*Ship, error) // Add this method
	Update(ship *Ship) error
	Delete(id uint) error
}
