package jwt

import (
	"errors"
	"eticket-api/config"
	"eticket-api/internal/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenUtil struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
	Rolename string `json:"rolename,omitempty"` // optional for refresh
	Username string `json:"username,omitempty"` // optional for refresh
	jwt.RegisteredClaims
}

// Constructor (call this in Run() or main)
func New(cfg *config.Configuration) *TokenUtil {
	return &TokenUtil{
		secretKey:          []byte(cfg.Auth.SecretKey),
		accessTokenExpiry:  cfg.Auth.AccessTokenExpiry,
		refreshTokenExpiry: cfg.Auth.RefreshTokenExpiry,
	}
}

func (tm *TokenUtil) GenerateAccessToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(tm.accessTokenExpiry)

	claims := &Claims{
		UserID:   user.ID,
		RoleID:   user.Role.ID,
		Rolename: user.Role.RoleName,
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

func (tm *TokenUtil) GenerateRefreshToken(user *entity.User) (string, error) {
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

func (tm *TokenUtil) ValidateToken(tokenString string) (*Claims, error) {
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
