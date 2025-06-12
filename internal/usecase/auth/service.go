package auth

import (
	"context"
	"errors"
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	Tx             *tx.TxManager
	AuthRepository *repository.AuthRepository
	UserRepository *repository.UserRepository
	TokenUtil      *jwt.TokenUtil
}

func NewAuthUsecase(
	tx *tx.TxManager,
	auth_repository *repository.AuthRepository,
	user_repository *repository.UserRepository,
	tm *jwt.TokenUtil,
) *AuthUsecase {
	return &AuthUsecase{
		Tx:             tx,
		AuthRepository: auth_repository,
		UserRepository: user_repository,
		TokenUtil:      tm,
	}
}

// Login authenticates a user and returns access and refresh tokens.
func (au *AuthUsecase) Login(ctx context.Context, request *model.UserLoginRequest) (string, string, error) {
	if request.Username == "" || request.Password == "" {
		return "", "", errors.New("username and password are required")
	}

	var accessToken, refreshToken string

	err := au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		user, repoErr := au.UserRepository.GetByUsername(tx, request.Username)
		if repoErr != nil {
			return fmt.Errorf("failed to retrieve user: %w", repoErr)
		}
		if user == nil {
			return errors.New("invalid credentials")
		}

		if !utils.CheckPasswordHash(request.Password, user.Password) {
			return errors.New("invalid credentials")
		}

		var err error
		accessToken, err = au.TokenUtil.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate access token: %w", err)
		}

		refreshToken, err = au.TokenUtil.GenerateRefreshToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %w", err)
		}

		// ✅ Validate and extract claims from the refresh token
		claims, err := au.TokenUtil.ValidateToken(refreshToken)
		if err != nil {
			return fmt.Errorf("failed to parse refresh token: %w", err)
		}

		refreshTokenEntity := &entity.RefreshToken{
			ID:        uuid.MustParse(claims.ID),
			UserID:    user.ID,
			Revoked:   false,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := au.AuthRepository.Create(tx, refreshTokenEntity); err != nil {
			return fmt.Errorf("failed to store refresh token: %w", err)
		}

		return nil // ✅ Only return error inside tx
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (au *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := au.TokenUtil.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	var newAccessToken string

	err = au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		// Check if refresh token exists and is valid
		session, err := au.AuthRepository.GetRefreshToken(tx, claims.ID)
		if err != nil {
			return fmt.Errorf("failed to get refresh token: %w", err)
		}
		if session.Revoked || session.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("refresh session invalid or expired")
		}

		// Get user associated with token
		user, err := au.UserRepository.GetByID(tx, claims.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// Generate new access token
		newAccessToken, err = au.TokenUtil.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate new access token: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (au *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	claims, err := au.TokenUtil.ValidateToken(refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}
	// Parse token ID (jti)
	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		return fmt.Errorf("invalid token ID: %w", err)
	}

	return au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		au.AuthRepository.RevokeRefreshTokenByID(tx, tokenID)
		return nil
	})
}
