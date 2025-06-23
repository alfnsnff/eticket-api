package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

type PasswordReset struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint      `gorm:"column:user_id"` // Assuming foreign key to User.ID
	Token     string    `gorm:"column:token"`
	Issued    bool      `gorm:"column:issued"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type AuthRepository interface {
	Create(db *gorm.DB, refreshToken *RefreshToken) error
	Count(db *gorm.DB) (int64, error)
	GetAllRefreshToken(db *gorm.DB) ([]*RefreshToken, error)
	GetRefreshToken(db *gorm.DB, id string) (*RefreshToken, error)
	RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error
	CreatePasswordReset(db *gorm.DB, pr *PasswordReset) error
	GetByToken(db *gorm.DB, token string) (*PasswordReset, error)
	MarkAsUsed(db *gorm.DB, token string) error
	DeleteExpired(db *gorm.DB) error
}

type AuthUsecase interface {
}
