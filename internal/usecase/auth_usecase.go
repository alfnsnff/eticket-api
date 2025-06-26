package usecase

import (
	"context"
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/templates"
	"eticket-api/internal/common/token"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/mapper"
	"eticket-api/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecase struct {
	DB             *gorm.DB // Assuming you have a DB field for the transaction manager
	AuthRepository domain.AuthRepository
	UserRepository domain.UserRepository
	Mailer         mailer.Mailer
	TokenUtil      token.TokenUtil
}

func NewAuthUsecase(
	db *gorm.DB,
	auth_repository domain.AuthRepository,
	user_repository domain.UserRepository,
	mailer mailer.Mailer,
	tm token.TokenUtil,
) *AuthUsecase {
	return &AuthUsecase{
		DB:             db,
		AuthRepository: auth_repository,
		UserRepository: user_repository,
		Mailer:         mailer,
		TokenUtil:      tm,
	}
}

// Login authenticates a user and returns access and refresh tokens.
func (au *AuthUsecase) Login(ctx context.Context, request *model.WriteLoginRequest) (*model.ReadLoginResponse, string, string, error) {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	user, err := au.UserRepository.FindByUsername(tx, request.Username)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(request.Password, user.Password) {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := au.TokenUtil.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := au.TokenUtil.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	claims, err := au.TokenUtil.ValidateToken(refreshToken)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to parse refresh token: %w", err)
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

	if err := au.AuthRepository.InsertRefreshToken(tx, refreshTokendomain); err != nil {
		return nil, "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, "", "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mapper.AuthToResponse(user), accessToken, refreshToken, nil
}

func (au *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	claims, err := au.TokenUtil.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists and is valid
	validSession, err := au.AuthRepository.FindRefreshTokenByIDAndStatus(tx, claims.ID, false)
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	if validSession.Revoked || validSession.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh session invalid or expired")
	}
	if validSession == nil {
		return "", errs.ErrNotFound
	}

	// Get user associated with token
	user, err := au.UserRepository.FindByID(tx, claims.User.ID)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return "", errs.ErrNotFound
	}

	// Generate new access token
	newAccessToken, err := au.TokenUtil.GenerateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newAccessToken, nil
}

func (au *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	claims, err := au.TokenUtil.ValidateToken(refreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}
	// Parse token ID (jti)
	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		return fmt.Errorf("invalid token ID: %w", err)
	}

	if err := au.AuthRepository.RevokeRefreshTokenByID(tx, tokenID); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (au *AuthUsecase) Me(ctx context.Context, accessToken string) (*model.ReadUserResponse, error) {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Parse and validate token
	claims, err := au.TokenUtil.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	user, err := au.UserRepository.FindByID(tx, claims.User.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, errs.ErrNotFound
	}

	return mapper.UserToResponse(user), nil
}

func (au *AuthUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	user, err := au.UserRepository.FindByEmail(tx, email)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return errs.ErrNotFound
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	reset := &domain.PasswordReset{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		CreatedAt: time.Now(),
	}

	if err := au.AuthRepository.InsertPasswordReset(tx, reset); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
	subject := "Password Reset"
	htmlBody := templates.PasswordResetEmail(user.Username, resetLink, time.Now().Year())

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	if err := au.Mailer.Send(user.Email, subject, htmlBody); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, token string, password string) error {
	tx := au.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	validReset, err := au.AuthRepository.FindPasswordResetByTokenAndStatus(tx, token, false)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token: %w", err)
	}
	if validReset.Issued || time.Now().After(validReset.ExpiresAt) {
		return errors.New("token expired or already used")
	}

	user, err := au.UserRepository.FindByID(tx, validReset.UserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return errs.ErrNotFound
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	if err := au.UserRepository.Update(tx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := au.AuthRepository.RevokePasswordResetByToken(tx, token); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}
