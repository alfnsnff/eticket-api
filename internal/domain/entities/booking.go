package entities

import "time"

type Booking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ScheduleID  uint      `gorm:"not null" json:"schedule_id"` // Foreign key
	CusName     string    `gorm:"not null" json:"cus_name"`
	PersonID    uint      `gorm:"not null" json:"person_id"`
	PhoneNumber uint      `gorm:"not null" json:"phone_number"`
	Email       string    `gorm:"not null" json:"email_address"`
	Birth       time.Time `gorm:"not null" json:"birth_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"` // Gorm will create the relationship
}

type BookingRepositoryInterface interface {
	Create(booking *Booking) error
	GetAll() ([]*Booking, error)
	GetByID(id uint) (*Booking, error) // Add this method
	Update(booking *Booking) error
	Delete(id uint) error
}
