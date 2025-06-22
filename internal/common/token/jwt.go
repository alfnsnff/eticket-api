package token

import (
	"errors"
	"eticket-api/config"
	constant "eticket-api/internal/common/constants"
	"eticket-api/internal/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenUtil interface {
	GenerateAccessToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateToken(token string) (*Claims, error)
}

type JWT struct {
	secretKey []byte
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
	Rolename string `json:"rolename,omitempty"` // optional for refresh
	Username string `json:"username,omitempty"` // optional for refresh
	jwt.RegisteredClaims
}

// Constructor (call this in Run() or main)
func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		secretKey: []byte(cfg.Token.SecretKey),
	}
}

func (tm *JWT) GenerateAccessToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(constant.AccessTokenExpiry)

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

func (tm *JWT) GenerateRefreshToken(user *entity.User) (string, error) {
	expirationTime := time.Now().Add(constant.RefreshTokenExpiry)

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

func (tm *JWT) ValidateToken(tokenString string) (*Claims, error) {
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
