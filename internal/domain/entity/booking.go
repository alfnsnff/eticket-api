package entity

import "time"

type Booking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ScheduleID   uint      `gorm:"not null;index;" json:"schedule_id"` // Foreign key
	PersonID     uint      `gorm:"not null" json:"person_id"`
	IdType       string    `gorm:"type:varchar(10);not null" json:"id_type"` // Changed to string to support leading zeros
	CustomerName string    `gorm:"not null" json:"customer_name"`
	PhoneNumber  string    `gorm:"type:varchar(15);not null" json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `gorm:"not null" json:"email"`
	Birth        time.Time `gorm:"not null" json:"birth_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
