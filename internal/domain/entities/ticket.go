package entities

import "time"

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookingID     uint      `gorm:"not null;index" json:"booking_id"` // Foreign key
	FareID        uint      `gorm:"not null;index" json:"fare_id"`    // Foreign key
	ScheduleID    uint      `gorm:"not null;index" json:"schedule_id"`
	PassengerName string    `json:"passenger_name"`
	SeatNumber    string    `json:"seat_number"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationship:
	Fare     Fare     `gorm:"foreignKey:FareID" json:"fare"`
	Booking  Booking  `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;" json:"booking"`
	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
}
