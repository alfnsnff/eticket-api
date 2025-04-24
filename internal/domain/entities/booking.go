package entities

import "time"

type Booking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ScheduleID  uint      `gorm:"not null" json:"schedule_id"` // Foreign key
	CusName     string    `gorm:"not null" json:"cus_name"`
	PersonID    uint      `gorm:"not null" json:"person_id"`
	PhoneNumber string    `gorm:"not null" json:"phone_number"` // Changed to string to support leading zeros
	Email       string    `gorm:"not null" json:"email_address"`
	Birth       time.Time `gorm:"not null" json:"birth_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
	Tickets  []Ticket `gorm:"foreignKey:BookingID" json:"tickets"` // ✅ Correct reference

}

// BookingRepositoryInterface defines repository methods
type BookingRepositoryInterface interface {
	Create(booking *Booking) error
	GetAll() ([]*Booking, error)
	GetByID(id uint) (*Booking, error)
	Update(booking *Booking) error
	Delete(id uint) error
}
