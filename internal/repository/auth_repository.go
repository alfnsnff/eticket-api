package repository

import (
	"errors"
	"eticket-api/internal/domain"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

// Standard CRUD operations for RefreshToken
func (ar *AuthRepository) CountRefreshToken(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.RefreshToken{}).Count(&total)
	return total, result.Error
}

func (ar *AuthRepository) InsertRefreshToken(db *gorm.DB, refreshToken *domain.RefreshToken) error {
	result := db.Create(refreshToken)
	return result.Error
}

func (ar *AuthRepository) InsertRefreshTokenBulk(db *gorm.DB, refreshToken []*domain.RefreshToken) error {
	result := db.Create(refreshToken)
	return result.Error
}

func (ar *AuthRepository) UpdateRefreshToken(db *gorm.DB, refreshToken *domain.RefreshToken) error {
	result := db.Save(refreshToken)
	return result.Error
}

func (ar *AuthRepository) UpdateRefreshTokenBulk(db *gorm.DB, refreshToken []*domain.RefreshToken) error {
	result := db.Save(refreshToken)
	return result.Error
}

func (ar *AuthRepository) DeleteRefreshToken(db *gorm.DB, id string) error {
	result := db.Delete(&domain.RefreshToken{}, "token = ?", id)
	return result.Error
}

func (ar *AuthRepository) FindAllRefreshToken(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.RefreshToken, error) {
	refreshToken := []*domain.RefreshToken{}
	query := db
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

func (ar *AuthRepository) FindRefreshTokenByID(db *gorm.DB, id string) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := db.First(&refreshToken, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (ar *AuthRepository) FindRefreshTokenByIDAndStatus(db *gorm.DB, id string, status bool) (*domain.RefreshToken, error) {
	refreshToken := new(domain.RefreshToken)
	result := db.First(&refreshToken, "id = ? AND revoked = ?", id, status)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return refreshToken, result.Error
}

func (aur *AuthRepository) RevokeRefreshTokenByID(db *gorm.DB, id uuid.UUID) error {
	return db.Model(&domain.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (ar *AuthRepository) CountPasswordResets(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.RefreshToken{}).Count(&total)
	return total, result.Error
}

// Password Reset operations
func (ar *AuthRepository) InsertPasswordReset(db *gorm.DB, passwordReset *domain.PasswordReset) error {
	result := db.Create(passwordReset)
	return result.Error
}

func (ar *AuthRepository) InsertPasswordResetBulk(db *gorm.DB, passwordResets []*domain.PasswordReset) error {
	result := db.Create(passwordResets)
	return result.Error
}

func (ar *AuthRepository) UpdatePasswordReset(db *gorm.DB, passwordReset *domain.PasswordReset) error {
	result := db.Save(passwordReset)
	return result.Error
}

func (ar *AuthRepository) UpdatePasswordResetBulk(db *gorm.DB, passwordResets []*domain.PasswordReset) error {
	result := db.Save(passwordResets)
	return result.Error
}

func (ar *AuthRepository) DeletePasswordReset(db *gorm.DB, id string) error {
	result := db.Delete(&domain.PasswordReset{}, "id = ?", id)
	return result.Error
}

func (ar *AuthRepository) FindAllPasswordResets(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.PasswordReset, error) {
	passwordResets := []*domain.PasswordReset{}
	query := db
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

func (ar *AuthRepository) FindPasswordResetByID(db *gorm.DB, id string) (*domain.PasswordReset, error) {
	passwordReset := new(domain.PasswordReset)
	result := db.First(&passwordReset, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return passwordReset, result.Error
}

func (ar *AuthRepository) FindPasswordResetByTokenAndStatus(db *gorm.DB, token string, status bool) (*domain.PasswordReset, error) {
	passwordReset := new(domain.PasswordReset)
	result := db.First(&passwordReset, "token = ? AND used = ?", token, status)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return passwordReset, result.Error
}

func (ar *AuthRepository) RevokePasswordResetByToken(db *gorm.DB, token string) error {
	passwordReset := new(domain.PasswordReset)
	result := db.Where("token = ?", token).
		First(passwordReset).
		Update("used", true)
	return result.Error
}
