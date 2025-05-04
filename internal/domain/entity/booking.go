package entity

import "time"

type Booking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ScheduleID   uint      `gorm:"not null;index;" json:"schedule_id"`         // Foreign key
	IDType       string    `gorm:"type:varchar(24);not null" json:"id_type"`   // Changed to string to support leading zeros
	IDNumber     string    `gorm:"type:varchar(24);not null" json:"id_number"` // Changed to string to support leading zeros
	CustomerName string    `gorm:"not null" json:"customer_name"`
	PhoneNumber  string    `gorm:"type:varchar(15);not null" json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `gorm:"not null" json:"email"`
	BirthDate    time.Time `gorm:"not null" json:"birth_date"`

	BookingTimestamp time.Time `gorm:"not null" json:"booking_timestamp"`       // Timestamp when the booking was confirmed
	TotalAmount      float32   `gorm:"not null" json:"total_amount"`            // Total price of all tickets in this booking
	Status           string    `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
