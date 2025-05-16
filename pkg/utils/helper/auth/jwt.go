package auth

import (
	"errors"
	"eticket-api/config"
	authentity "eticket-api/internal/domain/entity/auth"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var cfg = config.AppConfig.Auth

// Load the JWT signing key from environment variables (DO NOT hardcode in production)
var jwtSecretKey = []byte(cfg.SecretKey)

type Claims struct {
	UserID   uint   `json:"user_id"` // 👈 No longer shadows RegisteredClaims.ID
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed JWT access token (short-lived).
func GenerateAccessToken(user *authentity.User) (string, error) {
	expirationTime := time.Now().Add(cfg.AccessTokenExpiry * time.Minute)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eticket-api",              // Change to your app name
			Subject:   fmt.Sprintf("%d", user.ID), // User ID as subject
			ID:        uuid.New().String(),        // Unique ID for potential revocation tracking
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accesstoken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accesstoken, nil
}

// GenerateRefreshToken creates a signed JWT refresh token (long-lived).
func GenerateRefreshToken(user *authentity.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eticket-api",
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        uuid.New().String(), // Unique ID for potential revocation tracking
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshtoken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return refreshtoken, nil
}

// ValidateToken verifies a JWT token string and extracts the claims.
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
