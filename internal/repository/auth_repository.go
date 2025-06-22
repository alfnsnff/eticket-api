package repository

import (
	"eticket-api/internal/domain"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (aur *AuthRepository) Create(db *gorm.DB, refresh_token *domain.RefreshToken) error {
	result := db.Create(refresh_token)
	return result.Error
}

func (aur *AuthRepository) Count(db *gorm.DB) (int64, error) {
	refreshTokens := []*domain.RefreshToken{}
	var total int64
	result := db.Find(&refreshTokens).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (aur *AuthRepository) GetAllRefreshToken(db *gorm.DB) ([]*domain.RefreshToken, error) {
	tokens := []*domain.RefreshToken{}
	result := db.Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
}

func (aur *AuthRepository) GetRefreshToken(db *gorm.DB, id string) (*domain.RefreshToken, error) {
	token := new(domain.RefreshToken)
	result := db.First(&token, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (aur *AuthRepository) RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&domain.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

// Create a password reset token
func (aur *AuthRepository) CreatePasswordReset(db *gorm.DB, pr *domain.PasswordReset) error {
	return db.Create(pr).Error
}

// Get by token
func (aur *AuthRepository) GetByToken(db *gorm.DB, token string) (*domain.PasswordReset, error) {
	var pr domain.PasswordReset
	err := db.First(&pr, "token = ? AND used = false", token).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// Mark token as used
func (aur *AuthRepository) MarkAsUsed(db *gorm.DB, token string) error {
	return db.Model(&domain.PasswordReset{}).
		Where("token = ?", token).
		Update("used", true).Error
}

// Delete expired tokens
func (aur *AuthRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).
		Delete(&domain.PasswordReset{}).Error
}
