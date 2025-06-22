package repository

import (
	"eticket-api/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (aur *AuthRepository) Create(db *gorm.DB, refresh_token *entity.RefreshToken) error {
	result := db.Create(refresh_token)
	return result.Error
}

func (aur *AuthRepository) Count(db *gorm.DB) (int64, error) {
	refreshTokens := []*entity.RefreshToken{}
	var total int64
	result := db.Find(&refreshTokens).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (aur *AuthRepository) GetAllRefreshToken(db *gorm.DB) ([]*entity.RefreshToken, error) {
	tokens := []*entity.RefreshToken{}
	result := db.Find(&tokens)
	if result.Error != nil {
		return nil, result.Error
	}
	return tokens, nil
}

func (aur *AuthRepository) GetRefreshToken(db *gorm.DB, id string) (*entity.RefreshToken, error) {
	token := new(entity.RefreshToken)
	result := db.First(&token, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return token, nil
}

func (aur *AuthRepository) RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

// Create a password reset token
func (aur *AuthRepository) CreatePasswordReset(db *gorm.DB, pr *entity.PasswordReset) error {
	return db.Create(pr).Error
}

// Get by token
func (aur *AuthRepository) GetByToken(db *gorm.DB, token string) (*entity.PasswordReset, error) {
	var pr entity.PasswordReset
	err := db.First(&pr, "token = ? AND used = false", token).Error
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// Mark token as used
func (aur *AuthRepository) MarkAsUsed(db *gorm.DB, token string) error {
	return db.Model(&entity.PasswordReset{}).
		Where("token = ?", token).
		Update("used", true).Error
}

// Delete expired tokens
func (aur *AuthRepository) DeleteExpired(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).
		Delete(&entity.PasswordReset{}).Error
}
