package token

import (
	"errors"
	"eticket-api/config"
	constant "eticket-api/internal/common/constants"
	"eticket-api/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT struct {
	secretKey []byte
}

type Claims struct {
	User     *domain.User `json:"user,omitempty"`
	Rolename string       `json:"rolename,omitempty"`
	jwt.RegisteredClaims
}

// Constructor (call this in Run() or main)
func NewJWT(cfg *config.Config) *JWT {
	return &JWT{
		secretKey: []byte(cfg.Token.SecretKey),
	}
}

func (tm *JWT) GenerateAccessToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(constant.AccessTokenExpiry)

	claims := &Claims{
		User:     user,
		Rolename: user.Role.RoleName,
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

func (tm *JWT) GenerateRefreshToken(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(constant.RefreshTokenExpiry)

	claims := &Claims{
		User:     user,
		Rolename: user.Role.RoleName,
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
