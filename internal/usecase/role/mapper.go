package role

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Role entity to ReadRoleResponse model
func ToReadRoleResponse(role *entity.Role) *model.ReadRoleResponse {
	if role == nil {
		return nil
	}

	return &model.ReadRoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// Map a slice of Role entities to ReadRoleResponse models
func ToReadRoleResponses(roles []*entity.Role) []*model.ReadRoleResponse {
	responses := make([]*model.ReadRoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToReadRoleResponse(role)
	}
	return responses
}

// Map WriteRoleRequest model to Role entity
func FromWriteRoleRequest(request *model.WriteRoleRequest) *entity.Role {
	return &entity.Role{
		RoleName:    request.RoleName,
		Description: request.Description,
	}
}

// Map UpdateRoleRequest model to Role entity
func FromUpdateRoleRequest(request *model.UpdateRoleRequest, role *entity.Role) {
	role.RoleName = request.RoleName
	role.Description = request.Description
}
