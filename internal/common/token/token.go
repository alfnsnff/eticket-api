package token

import "eticket-api/internal/domain"

type TokenUtil interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken(user *domain.User) (string, error)
	ValidateToken(token string) (*Claims, error)
}
