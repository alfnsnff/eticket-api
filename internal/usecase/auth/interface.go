package auth

import (
	"eticket-api/internal/contracts"
)

// Use type aliases to point to the shared contracts
type (
	AuthRepository = contracts.AuthRepository
	UserRepository = contracts.UserRepository
)
