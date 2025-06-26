package model

import (
	"time"
)

// ClassDTO represents ticket class information.
type ReadClassResponse struct {
	ID         uint      `json:"id"`
	ClassName  string    `json:"class_name"`
	ClassAlias string    `json:"class_alias,omitempty"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WriteClassRequest struct {
	ClassName  string `json:"class_name" validate:"required"`   // Required
	Type       string `json:"type" validate:"required"`         // Adjust allowed values as needed
	ClassAlias string `json:"class_alias"  validate:"required"` // Optional
}

type UpdateClassRequest struct {
	ID         uint   `json:"id" validate:"required"`           // Required
	ClassName  string `json:"class_name" validate:"required"`   // Required
	Type       string `json:"type" validate:"required"`         // Adjust as needed
	ClassAlias string `json:"class_alias"  validate:"required"` // Optional
}
