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
	Role LoginUserRole `json:"role"`
}

type LoginUserRole struct {
	RoleName string `json:"role_name"`
}

func LoginFromRequest(request *LoginRequest) *domain.Login {
	return &domain.Login{
		Username: request.Username,
		Password: request.Password,
	}
}

func LoginToResponse(user *domain.User) *LoginResponse {
	return &LoginResponse{
		Role: LoginUserRole{
			RoleName: user.Role.RoleName,
		},
	}
}
