package request

import "time"

type CreateRoleRequest struct {
	RoleName    string `json:"role_name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateRoleRequest struct {
	ID          uint   `json:"id" validate:"required"`
	RoleName    string `json:"role_name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type RoleResponse struct {
	ID          uint      `json:"id"`
	RoleName    string    `json:"role_name"` // e.g., "admin", "editor"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
