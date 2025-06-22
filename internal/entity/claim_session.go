package entity

import (
	"time"
)

type ClaimSession struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID  string    `gorm:"column:session_id;type:uuid;unique;not null"`
	ScheduleID uint      `gorm:"column:schedule_id;not null;index"`
	Status     string    `gorm:"column:status;type:varchar(24);not null"` //
	ExpiresAt  time.Time `gorm:"column:expires_at;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`

	Schedule Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"` // Gorm will create the relationship
	Tickets  []Ticket `gorm:"foreignKey:ClaimSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tickets"`
}

func (cs *ClaimSession) TableName() string {
	return "claim_session"
}
