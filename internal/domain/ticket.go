package domain

import "time"

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ClassID       uint      `gorm:"not null" json:"class_id"` // Foreign key
	PassengerName string    `json:"passenger_name"`
	SeatNumber    string    `json:"seat_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationship:
	Class Class `gorm:"foreignKey:ClassID" json:"class"` // Gorm will create the relationship
}

type TicketRepository interface {
	Create(ticket *Ticket) error
	GetAll() ([]*Ticket, error)
	GetByID(id uint) (*Ticket, error) // Add this method
	Update(ticket *Ticket) error
	Delete(id uint) error
}
