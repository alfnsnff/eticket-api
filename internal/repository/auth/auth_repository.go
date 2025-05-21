package repository

import (
	entity "eticket-api/internal/domain/entity/auth"
	"eticket-api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct {
	repository.Repository[entity.RefreshToken]
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
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

func (aur *AuthRepository) RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&entity.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}
