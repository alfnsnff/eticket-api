package auth

import (
	"eticket-api/internal/domain"
)

// Use type aliases to point to the shared domain
type (
	AuthRepository = domain.AuthRepository
	UserRepository = domain.UserRepository
)
