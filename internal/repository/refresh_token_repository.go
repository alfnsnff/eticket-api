package repository

import (
	"context"
	"errors"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

// Standard CRUD operations for RefreshToken
func (r *RefreshTokenRepository) CountRefreshToken(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.RefreshToken{}).Count(&total)
	return total, result.Error
}

func (r *RefreshTokenRepository) InsertRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *domain.RefreshToken) error {
	result := conn.Create(refreshToken)
	return result.Error
}

func (r *RefreshTokenRepository) InsertRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshToken []*domain.RefreshToken) error {
	result := conn.Create(refreshToken)
	return result.Error
}

func (r *RefreshTokenRepository) UpdateRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *domain.RefreshToken) error {
	result := conn.Save(refreshToken)
	return result.Error
}

func (r *RefreshTokenRepository) UpdateRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshToken []*domain.RefreshToken) error {
	result := conn.Save(refreshToken)
	return result.Error
}

func (r *RefreshTokenRepository) DeleteRefreshToken(ctx context.Context, conn gotann.Connection, id string) error {
	result := conn.Delete(&domain.RefreshToken{}, "token = ?", id)
	return result.Error
}

func (r *RefreshTokenRepository) FindAllRefreshToken(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.RefreshToken, error) {
	refreshToken := []*domain.RefreshToken{}
	query := conn.Model(&domain.RefreshToken{})
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("user_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&refreshToken).Error
	return refreshToken, err
}

func (r *RefreshTokenRepository) FindRefreshTokenByID(ctx context.Context, conn gotann.Connection, id string) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := conn.First(&refreshToken, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (r *RefreshTokenRepository) FindRefreshTokenByIDAndStatus(ctx context.Context, conn gotann.Connection, id string, status bool) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := conn.First(&refreshToken, "id = ? AND revoked = ?", id, status)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (r *RefreshTokenRepository) RevokeRefreshTokenByID(ctx context.Context, conn gotann.Connection, id uuid.UUID) error {
	return conn.Model(&domain.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}
