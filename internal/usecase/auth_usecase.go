package usecase

import (
	"context"
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/transact"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthUsecase struct {
	Transactor     *transact.Transactor
	AuthRepository domain.AuthRepository
	UserRepository domain.UserRepository
	Mailer         mailer.Mailer
	TokenUtil      token.TokenUtil
}

func NewAuthUsecase(
	transactor *transact.Transactor,
	auth_repository domain.AuthRepository,
	user_repository domain.UserRepository,
	mailer mailer.Mailer,
	tm token.TokenUtil,
) *AuthUsecase {
	return &AuthUsecase{
		Transactor:     transactor,
		AuthRepository: auth_repository,
		UserRepository: user_repository,
		Mailer:         mailer,
		TokenUtil:      tm,
	}
}

func (uc *AuthUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.User, string, string, error) {
	var err error
	var user *domain.User
	var accessToken string
	var refreshToken string
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		user, err = uc.UserRepository.FindByUsername(ctx, tx, request.Username)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}
		if user == nil {
			return errors.New("invalid credentials")
		}
		if !utils.CheckPasswordHash(request.Password, user.Password) {
			return errors.New("invalid credentials")
		}
		accessToken, err = uc.TokenUtil.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate access token: %w", err)
		}
		refreshToken, err = uc.TokenUtil.GenerateRefreshToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate refresh token: %w", err)
		}
		claims, err := uc.TokenUtil.ValidateToken(refreshToken)
		if err != nil {
			return fmt.Errorf("failed to parse refresh token: %w", err)
		}

		refreshTokendomain := &domain.RefreshToken{
			ID:        uuid.MustParse(claims.ID),
			UserID:    user.ID,
			Revoked:   false,
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err = uc.AuthRepository.InsertRefreshToken(ctx, tx, refreshTokendomain); err != nil {
			if errs.IsUniqueConstraintError(err) {
				return errs.ErrConflict
			}
			return fmt.Errorf("failed to store refresh token: %w", err)
		}

		return nil
	}); err != nil {
		return nil, "", "", fmt.Errorf("login failed: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	var newAccessToken string
	if err := uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claims, err := uc.TokenUtil.ValidateToken(refreshToken)
		if err != nil {
			return fmt.Errorf("invalid refresh token: %w", err)
		}
		validSession, err := uc.AuthRepository.FindRefreshTokenByIDAndStatus(ctx, tx, claims.ID, false)
		if err != nil {
			return fmt.Errorf("failed to get refresh token: %w", err)
		}
		if validSession.Revoked || validSession.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("refresh session invalid or expired")
		}
		if validSession == nil {
			return errs.ErrNotFound
		}

		user, err := uc.UserRepository.FindByID(ctx, tx, claims.User.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}
		if user == nil {
			return errs.ErrNotFound
		}

		newAccessToken, err = uc.TokenUtil.GenerateAccessToken(user)
		if err != nil {
			return fmt.Errorf("failed to generate new access token: %w", err)
		}
		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return newAccessToken, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claims, err := uc.TokenUtil.ValidateToken(refreshToken)
		if err != nil {
			return fmt.Errorf("invalid refresh token: %w", err)
		}
		// Parse token ID (jti)
		tokenID, err := uuid.Parse(claims.ID)
		if err != nil {
			return fmt.Errorf("invalid token ID: %w", err)
		}

		if err := uc.AuthRepository.RevokeRefreshTokenByID(ctx, tx, tokenID); err != nil {
			return fmt.Errorf("failed to revoke refresh token: %w", err)
		}
		return nil
	})
}

func (uc *AuthUsecase) Me(ctx context.Context, accessToken string) (*domain.User, error) {
	var err error
	var claims *token.Claims
	var user *domain.User
	if err = uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
		claims, err = uc.TokenUtil.ValidateToken(accessToken)
		if err != nil {
			return fmt.Errorf("invalid access token: %w", err)
		}

		user, err = uc.UserRepository.FindByID(ctx, tx, claims.User.ID)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}
		if user == nil {
			return errs.ErrNotFound
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

// func (uc *AuthUsecase) RequestPasswordReset(ctx context.Context, email string) error {
// 	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
// 		user, err := uc.UserRepository.FindByEmail(ctx, tx, email)
// 		if err != nil {
// 			return fmt.Errorf("failed to retrieve user: %w", err)
// 		}
// 		if user == nil {
// 			return errs.ErrNotFound
// 		}

// 		token, err := utils.GenerateSecureToken(32)
// 		if err != nil {
// 			return fmt.Errorf("failed to generate token: %w", err)
// 		}

// 		reset := &domain.PasswordReset{
// 			UserID:    user.ID,
// 			Token:     token,
// 			ExpiresAt: time.Now().Add(15 * time.Minute),
// 			CreatedAt: time.Now(),
// 		}

// 		if err := uc.AuthRepository.InsertPasswordReset(ctx, tx, reset); err != nil {
// 			if errs.IsUniqueConstraintError(err) {
// 				return errs.ErrConflict
// 			}
// 			return fmt.Errorf("failed to save reset token: %w", err)
// 		}

// 		resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
// 		subject := "Password Reset"
// 		htmlBody := templates.PasswordResetEmail(user.Username, resetLink, time.Now().Year())

// 		if err := uc.Mailer.Send(user.Email, subject, htmlBody); err != nil {
// 			return fmt.Errorf("failed to send reset email: %w", err)
// 		}

// 		return nil
// 	})
// }

// func (uc *AuthUsecase) ResetPassword(ctx context.Context, token string, password string) error {
// 	return uc.Transactor.Execute(ctx, func(tx gotann.Transaction) error {
// 		validReset, err := uc.AuthRepository.FindPasswordResetByTokenAndStatus(ctx, tx, token, false)
// 		if err != nil {
// 			return fmt.Errorf("invalid or expired reset token: %w", err)
// 		}
// 		if validReset.Issued || time.Now().After(validReset.ExpiresAt) {
// 			return errors.New("token expired or already used")
// 		}

// 		user, err := uc.UserRepository.FindByID(ctx, tx, validReset.UserID)
// 		if err != nil {
// 			return fmt.Errorf("failed to retrieve user: %w", err)
// 		}
// 		if user == nil {
// 			return errs.ErrNotFound
// 		}

// 		hashedPassword, err := utils.HashPassword(password)
// 		if err != nil {
// 			return fmt.Errorf("failed to hash password: %w", err)
// 		}

// 		user.Password = hashedPassword

// 		if err := uc.UserRepository.Update(ctx, tx, user); err != nil {
// 			return fmt.Errorf("failed to update password: %w", err)
// 		}

// 		if err := uc.AuthRepository.RevokePasswordResetByToken(ctx, tx, token); err != nil {
// 			return fmt.Errorf("failed to mark token as used: %w", err)
// 		}
// 		return nil
// 	})
// }
