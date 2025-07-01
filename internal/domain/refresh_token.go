package domain

import (
	"context"
	"eticket-api/pkg/gotann"
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

type PasswordReset struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint      `gorm:"column:user_id"` // Assuming foreign key to User.ID
	Token     string    `gorm:"column:token"`
	Issued    bool      `gorm:"column:issued"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

type RefreshTokenRepository interface {
	CountRefreshToken(ctx context.Context, conn gotann.Connection) (int64, error)
	InsertRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *RefreshToken) error
	InsertRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshTokens []*RefreshToken) error
	UpdateRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *RefreshToken) error
	UpdateRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshTokens []*RefreshToken) error
	DeleteRefreshToken(ctx context.Context, conn gotann.Connection, id string) error
	FindAllRefreshToken(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*RefreshToken, error)
	FindRefreshTokenByID(ctx context.Context, conn gotann.Connection, id string) (*RefreshToken, error)
	FindRefreshTokenByIDAndStatus(ctx context.Context, conn gotann.Connection, id string, status bool) (*RefreshToken, error)
	RevokeRefreshTokenByID(ctx context.Context, conn gotann.Connection, id uuid.UUID) error
}
