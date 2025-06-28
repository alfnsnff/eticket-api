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

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

// Standard CRUD operations for RefreshToken
func (r *AuthRepository) CountRefreshToken(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.RefreshToken{}).Count(&total)
	return total, result.Error
}

func (r *AuthRepository) InsertRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *domain.RefreshToken) error {
	result := conn.Create(refreshToken)
	return result.Error
}

func (r *AuthRepository) InsertRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshToken []*domain.RefreshToken) error {
	result := conn.Create(refreshToken)
	return result.Error
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, conn gotann.Connection, refreshToken *domain.RefreshToken) error {
	result := conn.Save(refreshToken)
	return result.Error
}

func (r *AuthRepository) UpdateRefreshTokenBulk(ctx context.Context, conn gotann.Connection, refreshToken []*domain.RefreshToken) error {
	result := conn.Save(refreshToken)
	return result.Error
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, conn gotann.Connection, id string) error {
	result := conn.Delete(&domain.RefreshToken{}, "token = ?", id)
	return result.Error
}

func (r *AuthRepository) FindAllRefreshToken(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.RefreshToken, error) {
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

func (r *AuthRepository) FindRefreshTokenByID(ctx context.Context, conn gotann.Connection, id string) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := conn.First(&refreshToken, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (r *AuthRepository) FindRefreshTokenByIDAndStatus(ctx context.Context, conn gotann.Connection, id string, status bool) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := conn.First(&refreshToken, "id = ? AND revoked = ?", id, status)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (r *AuthRepository) RevokeRefreshTokenByID(ctx context.Context, conn gotann.Connection, id uuid.UUID) error {
	return conn.Model(&domain.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *AuthRepository) CountPasswordResets(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.RefreshToken{}).Count(&total)
	return total, result.Error
}

// Password Reset operations
func (r *AuthRepository) InsertPasswordReset(ctx context.Context, conn gotann.Connection, passwordReset *domain.PasswordReset) error {
	result := conn.Create(passwordReset)
	return result.Error
}

func (r *AuthRepository) InsertPasswordResetBulk(ctx context.Context, conn gotann.Connection, passwordResets []*domain.PasswordReset) error {
	result := conn.Create(passwordResets)
	return result.Error
}

func (r *AuthRepository) UpdatePasswordReset(ctx context.Context, conn gotann.Connection, passwordReset *domain.PasswordReset) error {
	result := conn.Save(passwordReset)
	return result.Error
}

func (r *AuthRepository) UpdatePasswordResetBulk(ctx context.Context, conn gotann.Connection, passwordResets []*domain.PasswordReset) error {
	result := conn.Save(passwordResets)
	return result.Error
}

func (r *AuthRepository) DeletePasswordReset(ctx context.Context, conn gotann.Connection, id string) error {
	result := conn.Delete(&domain.PasswordReset{}, "id = ?", id)
	return result.Error
}

func (r *AuthRepository) FindAllPasswordResets(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.PasswordReset, error) {
	passwordResets := []*domain.PasswordReset{}
	query := conn.Model(&domain.PasswordReset{})
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("user_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&passwordResets).Error
	return passwordResets, err
}

func (r *AuthRepository) FindPasswordResetByID(ctx context.Context, conn gotann.Connection, id string) (*domain.PasswordReset, error) {
	passwordReset := new(domain.PasswordReset)
	result := conn.First(&passwordReset, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return passwordReset, result.Error
}

func (r *AuthRepository) FindPasswordResetByTokenAndStatus(ctx context.Context, conn gotann.Connection, token string, status bool) (*domain.PasswordReset, error) {
	passwordReset := new(domain.PasswordReset)
	result := conn.First(&passwordReset, "token = ? AND used = ?", token, status)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return passwordReset, result.Error
}

func (r *AuthRepository) RevokePasswordResetByToken(zctx context.Context, conn gotann.Connection, token string) error {
	passwordReset := new(domain.PasswordReset)
	result := conn.Where("token = ?", token).
		First(passwordReset).
		Update("used", true)
	return result.Error
}
