package entities

import "time"

type Class struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RouteID   uint      `gorm:"not null" json:"route_id"` // Foreign key
	ClassName string    `gorm:"not null" json:"class_name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationship:
	// Tickets []Ticket `gorm:"foreignKey:ClassID" json:"tickets"`
	Route Route `gorm:"foreignKey:RouteID" json:"route"` // Gorm will create the relationship
}

type ClassRepositoryInterface interface {
	Create(class *Class) error
	GetAll() ([]*Class, error)
	GetByID(id uint) (*Class, error) // Add this method
	Update(class *Class) error
	Delete(id uint) error
}
