package entity

import "time"

type Allocation struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"not null" json:"schedule_id"` // Foreign key
	ClassID    uint      `gorm:"not null" json:"class_id"`    // Foreign key
	Quota      int       `gorm:"not null" json:"quota"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
