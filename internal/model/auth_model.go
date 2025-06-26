package model

type ReadLoginResponse struct {
	Role UserRole `json:"role"`
}

type WriteForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"` // Must be a valid email
}

type WriteLoginRequest struct {
	Username string `json:"username" validate:"required"` // Required field
	Password string `json:"password" validate:"required"` // Required field
}
