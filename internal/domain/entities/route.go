package entities

import "time"

type Route struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	DepartureHarborID uint      `json:"deaprature_harbor_id"`
	ArrivalHarborID   uint      `json:"arrival_harbor_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	DepartureHarbor Harbor `gorm:"foreignKey:DepartureHarborID" json:"departure_harbor"` // Gorm will create the relationship
	ArrivalHarbor   Harbor `gorm:"foreignKey:ArrivalHarborID" json:"arrival_harbor"`     // Gorm will create the relationship
}

type RouteRepositoryInterface interface {
	Create(route *Route) error
	GetAll() ([]*Route, error)
	GetByID(id uint) (*Route, error) // Add this method
	Update(route *Route) error
	Delete(id uint) error
}
