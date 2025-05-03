package entity

import "time"

type Ticket struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ScheduleID    uint      `gorm:"index;not null" json:"schedule_id"` // Corrected: Non-nullable uint
	ClassID       uint      `gorm:"index;not null" json:"class_id"`    // Corrected: Non-nullable uint
	Status        string    `gorm:"not null" json:"status"`            // Added not null, status is essential
	BookingID     *uint     `gorm:"index" json:"booking_id"`           // Retained, assuming Booking entity has a purpose (e.g., final transaction group)
	PassengerName *string   `json:"passenger_name"`
	SeatNumber    *string   `json:"seat_number"`                // Nullable is fine
	Price         float32   `gorm:"not null" json:"price"`      // Added not null, price is essential post-claim
	ExpiresAt     time.Time `gorm:"not null" json:"expires_at"` // Added not null, essential for timeout
	ClaimedAt     time.Time `gorm:"not null" json:"claimed_at"` // Added not null, essential for timeout
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
