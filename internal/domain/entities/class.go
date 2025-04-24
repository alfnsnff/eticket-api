package entities

import "time"

type Class struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClassRepositoryInterface interface {
	Create(class *Class) error
	GetAll() ([]*Class, error)
	GetByID(id uint) (*Class, error) // Add this method
	Update(class *Class) error
	Delete(id uint) error
}
