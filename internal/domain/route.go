package domain

import "time"

type Route struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	DepartureHarbor string    `json:"deaprature_harbor"`
	ArrivalHarbor   string    `json:"arrival_harbor"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RouteRepository interface {
	Create(route *Route) error
	GetAll() ([]*Route, error)
	GetByID(id uint) (*Route, error) // Add this method
	Update(route *Route) error
	Delete(id uint) error
}
