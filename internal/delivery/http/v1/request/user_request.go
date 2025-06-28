package request

import "time"

type CreateUserRequest struct {
	RoleID   uint   `json:"role_id" validate:"required"`
	Username string `json:"username" validate:"required,min=5"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type UpdateUserRequest struct {
	ID       uint   `json:"id" validate:"required"`
	RoleID   uint   `json:"role_id" validate:"required"`
	Username string `json:"username" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	FullName string `json:"full_name" validate:"required"`
}

type ReadUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRole struct {
	ID          uint   `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}
