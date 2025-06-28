package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

func AuthToResponse(user *domain.User) *model.ReadLoginResponse {
	var role model.UserRole

	if user.Role.ID != 0 {
		role = model.UserRole{
			ID:          user.Role.ID,
			RoleName:    user.Role.RoleName,
			Description: user.Role.Description,
		}
	} else {
		role = model.UserRole{
			ID:          0,
			RoleName:    "UNKNOWN",
			Description: "Unknown role",
		}
	}

	return &model.ReadLoginResponse{
		Role: role,
	}
}
