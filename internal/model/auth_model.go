package model

type UserRole struct {
	ID          uint   `json:"id"`
	RoleName    string `json:"role_name"` // e.g., "admin", "editor"
	Description string `json:"description"`
}

type ReadLoginResponse struct {
	Username string   `json:"username"` // e.g., "admin", "editor"
	Role     UserRole `json:"role"`
}

type WriteForgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"` // Must be a valid email
}

type WriteLoginRequest struct {
	Username string `json:"username" validate:"required"` // Required field
	Password string `json:"password" validate:"required"` // Required field
}
