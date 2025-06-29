package requests

import (
	"eticket-api/internal/domain"
	"time"
)

type CreateClassRequest struct {
	ClassName  string `json:"class_name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	ClassAlias string `json:"class_alias"  validate:"required"`
}

type UpdateClassRequest struct {
	ID         uint   `json:"id" validate:"required"`
	ClassName  string `json:"class_name" validate:"required"`
	Type       string `json:"type" validate:"required"`
	ClassAlias string `json:"class_alias"  validate:"required"`
}

type ClassResponse struct {
	ID         uint      `json:"id"`
	ClassName  string    `json:"class_name"`
	ClassAlias string    `json:"class_alias,omitempty"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Map Class domain to ReadClassResponse model
func ClassToResponse(class *domain.Class) *ClassResponse {
	return &ClassResponse{
		ID:         class.ID,
		ClassName:  class.ClassName,
		ClassAlias: class.ClassAlias,
		Type:       class.Type,
		CreatedAt:  class.CreatedAt,
		UpdatedAt:  class.UpdatedAt,
	}
}

func ClassFromCreate(request *CreateClassRequest) *domain.Class {
	return &domain.Class{
		ClassName:  request.ClassName,
		ClassAlias: request.ClassAlias,
		Type:       request.Type,
	}
}

func ClassFromUpdate(request *UpdateClassRequest) *domain.Class {
	return &domain.Class{
		ID:         request.ID,
		ClassName:  request.ClassName,
		ClassAlias: request.ClassAlias,
		Type:       request.Type,
	}
}
