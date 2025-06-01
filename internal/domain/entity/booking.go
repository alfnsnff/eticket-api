package entity

import "time"

type Booking struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	ScheduleID   uint      `gorm:"column:schedule_id;not null;index;"`
	IDType       string    `gorm:"column:id_type;type:varchar(24);not null"`
	IDNumber     string    `gorm:"column:id_number;type:varchar(24);not null"`
	CustomerName string    `gorm:"column:customer_name;type:varchar(32);not null"`
	PhoneNumber  string    `gorm:"column:phone_number;type:varchar(14);not null"`
	Email        string    `gorm:"column:email;not null"`
	Status       string    `gorm:"column:status;type:varchar(24);not null"`
	BookedAt     time.Time `gorm:"column:booked_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

func (b *Booking) TableName() string {
	return "booking"
}
