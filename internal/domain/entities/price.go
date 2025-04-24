package entities

import "time"

type Price struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RouteID     uint      `json:"route_id"`
	ShipClassID uint      `json:"ship_class_id"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	ShipClass ShipClass `gorm:"foreignKey:ShipClassID" json:"ship_class"`
}

type PriceRepositoryInterface interface {
	Create(price *Price) error
	GetAll() ([]*Price, error)
	GetByID(id uint) (*Price, error) // Add this method
	GetByIDs(priceIDs []uint) ([]*Price, error)
	GetByRouteID(routeID uint) ([]*Price, error)
	Update(price *Price) error
	Delete(id uint) error
}
