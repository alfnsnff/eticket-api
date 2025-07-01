package requests

import (
	"eticket-api/internal/domain"
	"time"
)

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

type UserResponse struct {
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

// Map User domain to ReadUserResponse model
func UserToResponse(user *domain.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role: UserRole{
			ID:          user.Role.ID,
			RoleName:    user.Role.RoleName,
			Description: user.Role.Description,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Map WriteUserRequest model to User domain
func UserFromCreate(request *CreateUserRequest) *domain.User {
	return &domain.User{
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
		Password: request.Password,
		RoleID:   request.RoleID,
	}
}

// Map WriteUserRequest model to User domain
func UserFromUpdate(request *UpdateUserRequest) *domain.User {
	return &domain.User{
		ID:       request.ID,
		Username: request.Username,
		Email:    request.Email,
		FullName: request.FullName,
		Password: request.Password,
		RoleID:   request.RoleID,
	}
}
