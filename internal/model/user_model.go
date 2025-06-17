package model

import "time"

type UserRole struct {
	ID          uint   `json:"id"`
	RoleName    string `json:"role_name"` // e.g., "admin", "editor"
	Description string `json:"description"`
}

type ReadUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"` // e.g., "admin", "editor"
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteUserRequest struct {
	RoleID   uint   `json:"role_id" binding:"required"`
	Username string `json:"username" binding:"required"` // e.g., "admin", "editor"
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

// AllocationDTO represents a Allocation.
type UpdateUserRequest struct {
	ID       uint   `json:"id"`
	RoleID   uint   `json:"role_id" binding:"required"`
	Username string `json:"username"` // e.g., "admin", "editor"
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// AllocationDTO represents a Allocation.
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"` // e.g., "admin", "editor"
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Optional: Include refresh token
}
