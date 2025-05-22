package jwt

import (
	"errors"
	"eticket-api/config"
	authentity "eticket-api/internal/domain/entity/auth"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username,omitempty"` // optional for refresh
	jwt.RegisteredClaims
}

// Constructor (call this in Run() or main)
func New(cfg *config.Config) *TokenManager {
	return &TokenManager{
		secretKey:          []byte(cfg.Auth.SecretKey),
		accessTokenExpiry:  cfg.Auth.AccessTokenExpiry,
		refreshTokenExpiry: cfg.Auth.RefreshTokenExpiry,
	}
}

func (tm *TokenManager) GenerateAccessToken(user *authentity.User) (string, error) {
	expirationTime := time.Now().Add(tm.accessTokenExpiry)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eticket-api",
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secretKey)
}

func (tm *TokenManager) GenerateRefreshToken(user *authentity.User) (string, error) {
	expirationTime := time.Now().Add(tm.refreshTokenExpiry)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eticket-api",
			Subject:   fmt.Sprintf("%d", user.ID),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secretKey)
}

func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
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

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
