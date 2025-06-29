package requests

import "eticket-api/internal/domain"

type ForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"` // Must be a valid email
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Required field
	Password string `json:"password" validate:"required"` // Required field
}

type LoginResponse struct {
	Role UserRole `json:"role"`
}

func LoginFromRequest(request *LoginRequest) *domain.LoginRequest {
	return &domain.LoginRequest{
		Username: request.Username,
		Password: request.Password,
	}
}
