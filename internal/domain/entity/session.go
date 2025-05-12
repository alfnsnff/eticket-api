package entity

import "time"

type ClaimSession struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID string    `gorm:"column:session_id;type:uuid;unique;not null"`
	ClaimedAt time.Time `gorm:"column:claimed_at;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (cs *ClaimSession) TableName() string {
	return "claim_session"
}
