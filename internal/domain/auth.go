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
	// RefreshToken CRUD operations
	CountRefreshToken(db *gorm.DB) (int64, error)
	InsertRefreshToken(db *gorm.DB, refreshToken *RefreshToken) error
	InsertRefreshTokenBulk(db *gorm.DB, refreshTokens []*RefreshToken) error
	UpdateRefreshToken(db *gorm.DB, refreshToken *RefreshToken) error
	UpdateRefreshTokenBulk(db *gorm.DB, refreshTokens []*RefreshToken) error
	DeleteRefreshToken(db *gorm.DB, id string) error
	FindAllRefreshToken(db *gorm.DB, limit, offset int, sort, search string) ([]*RefreshToken, error)
	FindRefreshTokenByID(db *gorm.DB, id string) (*RefreshToken, error)
	FindRefreshTokenByIDAndStatus(db *gorm.DB, id string, status bool) (*RefreshToken, error)
	RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error

	// PasswordReset CRUD operations
	CountPasswordResets(db *gorm.DB) (int64, error)
	InsertPasswordReset(db *gorm.DB, passwordReset *PasswordReset) error
	InsertPasswordResetBulk(db *gorm.DB, passwordResets []*PasswordReset) error
	UpdatePasswordReset(db *gorm.DB, passwordReset *PasswordReset) error
	UpdatePasswordResetBulk(db *gorm.DB, passwordResets []*PasswordReset) error
	DeletePasswordReset(db *gorm.DB, id string) error
	FindAllPasswordResets(db *gorm.DB, limit, offset int, sort, search string) ([]*PasswordReset, error)
	FindPasswordResetByID(db *gorm.DB, id string) (*PasswordReset, error) // Note: This should probably return *PasswordReset
	FindPasswordResetByTokenAndStatus(db *gorm.DB, token string, status bool) (*PasswordReset, error)
	RevokePasswordResetByToken(db *gorm.DB, token string) error
}

type AuthUsecase interface {
}
