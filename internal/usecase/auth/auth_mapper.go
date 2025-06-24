package auth

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map User domain to ReadUserResponse model
func UserToResponse(user *domain.User) *model.ReadUserResponse {
	return &model.ReadUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role: model.UserRole{
			ID:       user.Role.ID,
			RoleName: user.Role.RoleName,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Map User domain to ReadUserResponse model
func AuthToResponse(user *domain.User) *model.ReadLoginResponse {
	return &model.ReadLoginResponse{
		Username: user.Username,
		Role: model.UserRole{
			ID:       user.Role.ID,
			RoleName: user.Role.RoleName,
		},
	}
}
