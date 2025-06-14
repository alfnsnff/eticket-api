package auth

import (
	"context"
	"errors"
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/mailer"
	"eticket-api/internal/common/tx"
	"eticket-api/internal/common/utils"
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"eticket-api/internal/model/mapper"
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
	Mailer         *mailer.SMTPMailer
	TokenUtil      *jwt.TokenUtil
}

func NewAuthUsecase(
	tx *tx.TxManager,
	auth_repository *repository.AuthRepository,
	user_repository *repository.UserRepository,
	mailer *mailer.SMTPMailer,
	tm *jwt.TokenUtil,
) *AuthUsecase {
	return &AuthUsecase{
		Tx:             tx,
		AuthRepository: auth_repository,
		UserRepository: user_repository,
		Mailer:         mailer,
		TokenUtil:      tm,
	}
}

// Login authenticates a user and returns access and refresh tokens.
func (au *AuthUsecase) Login(ctx context.Context, request *model.UserLoginRequest) (*model.ReadUserResponse, string, string, error) {
	if request.Username == "" || request.Password == "" {
		return nil, "", "", errors.New("username and password are required")
	}

	var accessToken, refreshToken string
	var userd *entity.User

	err := au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		user, err := au.UserRepository.GetByUsername(tx, request.Username)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}
		if user == nil {
			return errors.New("invalid credentials")
		}
		userd = user

		if !utils.CheckPasswordHash(request.Password, user.Password) {
			return errors.New("invalid credentials")
		}

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
		return nil, "", "", err
	}

	fmt.Println(userd)
	return mapper.UserMapper.ToModel(userd), accessToken, refreshToken, nil
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

func (au *AuthUsecase) Me(ctx context.Context, accessToken string) (*model.ReadUserResponse, error) {
	// Parse and validate token
	claims, err := au.TokenUtil.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	var user *entity.User
	err = au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		u, err := au.UserRepository.GetByID(tx, claims.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		user = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return mapper.UserMapper.ToModel(user), nil
}

func (au *AuthUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	return au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		user, err := au.UserRepository.GetByEmail(tx, email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// don't reveal email existence
				return nil
			}
			return fmt.Errorf("failed to retrieve user: %w", err)
		}

		token, err := utils.GenerateSecureToken(32)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		reset := &entity.PasswordReset{
			UserID:    user.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now(),
		}

		if err := au.AuthRepository.CreatePasswordReset(tx, reset); err != nil {
			return fmt.Errorf("failed to save reset token: %w", err)
		}

		resetLink := fmt.Sprintf("https://yourdomain.com/reset-password?token=%s", token)
		subject := "Password Reset"
		htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="UTF-8">
		<title>Password Reset</title>
		<style>
			body {
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
			background-color: #f4f4f4;
			margin: 0;
			padding: 0;
			}
			.container {
			max-width: 600px;
			margin: 40px auto;
			background-color: #ffffff;
			padding: 30px;
			border-radius: 8px;
			box-shadow: 0 0 10px rgba(0,0,0,0.05);
			}
			.header {
			text-align: center;
			padding-bottom: 20px;
			}
			.header img {
			width: 120px;
			}
			.title {
			font-size: 24px;
			color: #333333;
			margin-bottom: 10px;
			}
			.content {
			font-size: 16px;
			color: #555555;
			line-height: 1.5;
			margin-bottom: 30px;
			}
			.button {
			display: inline-block;
			padding: 12px 20px;
			background-color: #6441a5;
			color: #ffffff;
			text-decoration: none;
			border-radius: 6px;
			font-weight: bold;
			}
			.footer {
			font-size: 13px;
			color: #888888;
			text-align: center;
			margin-top: 30px;
			}
		</style>
		</head>
		<body>
		<div class="container">
			<div class="header">
			<img src="https://upload.wikimedia.org/wikipedia/commons/1/13/Ticket_emoji.png" alt="Logo" />
			</div>
			<div class="title">Reset Your Password</div>
			<div class="content">
			Hello %s,<br><br>
			We received a request to reset your password. If you made this request, click the button below. This link will expire in 15 minutes.<br><br>
			<a href="%s" class="button">Reset Password</a><br><br>
			If you did not request a password reset, please ignore this email or contact support if you have concerns.
			</div>
			<div class="footer">
			&copy; %d eTicket. All rights reserved.
			</div>
		</div>
		</body>
		</html>
		`, user.Username, resetLink, time.Now().Year())

		return au.Mailer.Send(user.Email, subject, htmlBody)
	})
}

func (au *AuthUsecase) ResetPassword(ctx context.Context, token string, password string) error {
	return au.Tx.Execute(ctx, func(tx *gorm.DB) error {
		reset, err := au.AuthRepository.GetByToken(tx, token)
		if err != nil {
			return fmt.Errorf("invalid or expired reset token: %w", err)
		}
		if reset.Issued || time.Now().After(reset.ExpiresAt) {
			return errors.New("token expired or already used")
		}

		user, err := au.UserRepository.GetByID(tx, reset.UserID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		if err := au.UserRepository.UpdatePassword(tx, user.ID, hashedPassword); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		if err := au.AuthRepository.MarkAsUsed(tx, token); err != nil {
			return fmt.Errorf("failed to mark token as used: %w", err)
		}

		return nil
	})
}
