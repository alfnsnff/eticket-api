package user

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map User entity to ReadUserResponse model
func ToReadUserResponse(user *entity.User) *model.ReadUserResponse {
	if user == nil {
		return nil
	}

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

// Map a slice of User entities to ReadUserResponse models
func ToReadUserResponses(users []*entity.User) []*model.ReadUserResponse {
	responses := make([]*model.ReadUserResponse, len(users))
	for i, user := range users {
		responses[i] = ToReadUserResponse(user)
	}
	return responses
}

// Map WriteUserRequest model to User entity
func FromWriteUserRequest(request *model.WriteUserRequest, hashedPassword string) *entity.User {
	return &entity.User{
		RoleID:   request.RoleID,
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		FullName: request.FullName,
	}
}

// Map UpdateUserRequest model to User entity
func FromUpdateUserRequest(request *model.UpdateUserRequest, user *entity.User) {
	user.RoleID = request.RoleID
	user.Username = request.Username
	user.Email = request.Email
	user.FullName = request.FullName
	if request.Password != "" {
		user.Password = request.Password // Ensure password is hashed before updating
	}
}
