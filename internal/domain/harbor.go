package domain

import "time"

type Harbor struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	HarborName string    `gorm:"not null" json:"harbor_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type HarborRepositoryInterface interface {
	Create(harbor *Harbor) error
	GetAll() ([]*Harbor, error)
	GetByID(id uint) (*Harbor, error) // Add this method
	Update(harbor *Harbor) error
	Delete(id uint) error
}
