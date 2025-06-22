package role

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Role domain to ReadRoleResponse model
func RoleToResponse(role *domain.Role) *model.ReadRoleResponse {
	return &model.ReadRoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
