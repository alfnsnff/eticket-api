package request

import (
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
