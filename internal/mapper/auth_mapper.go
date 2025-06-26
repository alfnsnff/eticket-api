package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map User domain to ReadUserResponse model
func AuthToResponse(user *domain.User) *model.ReadLoginResponse {
	return &model.ReadLoginResponse{
		Role: model.UserRole{
			ID:          user.Role.ID,
			RoleName:    user.Role.RoleName,
			Description: user.Role.Description,
		},
	}
}
