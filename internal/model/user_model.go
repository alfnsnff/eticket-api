package model

import "time"

type ReadUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"` // e.g., "admin", "editor"
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteUserRequest struct {
	RoleID   uint   `json:"role_id" validate:"required"`
	Username string `json:"username" validate:"required,min=3"` // Optional: enforce a minimum length
	Email    string `json:"email" validate:"required,email"`    // Email format validation
	Password string `json:"password" validate:"required,min=6"` // Optional: minimum password length
	FullName string `json:"full_name" validate:"required"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id" validate:"required"` // ID must be provided for updates
	RoleID   uint   `json:"role_id" validate:"required"`
	Username string `json:"username" validate:"required,min=3"` // Still required, even if updating
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=6"` // Optional, but if filled must be valid
	FullName string `json:"full_name" validate:"required"`
}
