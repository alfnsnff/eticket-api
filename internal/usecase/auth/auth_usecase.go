package usecase

import (
	"context"
	"errors"
	authentity "eticket-api/internal/domain/entity/auth"
	authmodel "eticket-api/internal/model/auth"
	authrepository "eticket-api/internal/repository/auth"
	"eticket-api/pkg/jwt"
	utils "eticket-api/pkg/utils/hash"
	"eticket-api/pkg/utils/tx"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	Tx             *tx.TxManager
	AuthRepository *authrepository.AuthRepository
	UserRepository *authrepository.UserRepository
	TokenManager   *jwt.TokenManager
}

func NewAuthUsecase(
	tx *tx.TxManager,
	auth_repository *authrepository.AuthRepository,
	user_repository *authrepository.UserRepository,
	tm *jwt.TokenManager,
) *AuthUsecase {
	return &AuthUsecase{
		Tx:             tx,
		AuthRepository: auth_repository,
		UserRepository: user_repository,
		TokenManager:   tm,
	}
}

// Login authenticates a user and returns access and refresh tokens.
func (uc *AuthUsecase) Login(ctx context.Context, request *authmodel.UserLoginRequest) (string, string, error) {
	if request.Username == "" || request.Password == "" {
		return "", "", errors.New("username and password are required")
	}

	var accessToken, refreshToken string

	err := uc.Tx.Execute(ctx, func(txDB *gorm.DB) error {
		user, repoErr := uc.UserRepository.GetByUsername(txDB, request.Username)
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
		accessToken, err = uc.TokenManager.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate access token: %w", err)
		}

		refreshToken, err = uc.TokenManager.GenerateRefreshToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %w", err)
		}

		// ✅ Validate and extract claims from the refresh token
		claims, err := uc.TokenManager.ValidateToken(refreshToken)
		if err != nil {
			return fmt.Errorf("failed to parse refresh token: %w", err)
		}

		refreshTokenEntity := &authentity.RefreshToken{
			ID:        uuid.MustParse(claims.ID),
			UserID:    user.ID,
			Revoked:   false,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := uc.AuthRepository.Create(txDB, refreshTokenEntity); err != nil {
			return fmt.Errorf("failed to store refresh token: %w", err)
		}

		return nil // ✅ Only return error inside tx
	})

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// func (uc *AuthUsecase) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
// 	return uc.AuthRepository.RevokeRefreshTokenByID(, tokenID)
// }

func (au *AuthUsecase) RevokeRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	return au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		au.AuthRepository.RevokeRefreshTokenByID(tx, tokenID)
		return nil
	})
}
