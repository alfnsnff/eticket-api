package entities

import "time"

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookingID     uint      `gorm:"not null;index" json:"booking_id"` // Foreign key
	PriceID       uint      `gorm:"not null;index" json:"price_id"`   // Foreign key
	ScheduleID    uint      `gorm:"not null;index" json:"schedule_id"`
	PassengerName string    `json:"passenger_name"`
	SeatNumber    string    `json:"seat_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationship:
	Price    Price    `gorm:"foreignKey:PriceID" json:"price"`
	Booking  Booking  `gorm:"foreignKey:BookingID" json:"booking"` // Gorm will create the relationship
	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
}

type TicketRepositoryInterface interface {
	Create(ticket *Ticket) error
	GetAll() ([]*Ticket, error)
	GetByID(id uint) (*Ticket, error) // Add this method
	GetBookedCount(scheduleID uint, priceID uint) (int, error)
	Update(ticket *Ticket) error
	Delete(id uint) error
}
