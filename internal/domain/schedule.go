package domain

import "time"

type Schedule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RouteID   uint      `gorm:"not null" json:"route_id"` // Foreign key
	ShipID    uint      `gorm:"not null" json:"ship_id"`
	Datetime  time.Time `gorm:"not null" json:"datetime"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Route Route `gorm:"foreignKey:RouteID" json:"route"` // Gorm will create the relationship
	Ship  Ship  `gorm:"foreignKey:ShipID" json:"ship"`   // Gorm will create the relationship
}

type ScheduleRepositoryInterface interface {
	Create(schedule *Schedule) error
	GetAll() ([]*Schedule, error)
	GetByID(id uint) (*Schedule, error) // Add this method
	Update(schedule *Schedule) error
	Delete(id uint) error
}
