package request

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
