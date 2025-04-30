package entities

import "time"

type TicketReservation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"not null" json:"schedule_id"`
	SessionID  string    `gorm:"index;not null" json:"session_id"` // UUID from client or backend
	Status     string    `gorm:"not null" json:"status"`           // "held", "expired", "confirmed"
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
}
