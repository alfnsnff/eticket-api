package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"` // Correct type for UUID
	UserID    uint      `gorm:"column:user_id"`       // Assuming foreign key to User.ID
	Revoked   bool      `gorm:"column:revoked"`
	IssuedAt  time.Time `gorm:"column:issued_at"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
