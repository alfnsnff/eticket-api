package requests

import (
	"eticket-api/internal/domain"
	"time"
)

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

// Map Role RoleResponse model
func RoleToResponse(role *domain.Role) *RoleResponse {
	return &RoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func RoleFromCreate(request *CreateRoleRequest) *domain.Role {
	return &domain.Role{
		RoleName:    request.RoleName,
		Description: request.Description,
	}
}

func RoleFromUpdate(request *UpdateRoleRequest) *domain.Role {
	return &domain.Role{
		ID:          request.ID,
		RoleName:    request.RoleName,
		Description: request.Description,
	}
}
