package model

import (
	"time"
)

// ClassDTO represents ticket class information.
type ReadClassResponse struct {
	ID         uint      `json:"id"`
	ClassName  string    `json:"class_name"`
	ClassAlias *string   `json:"class_alias,omitempty"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type WriteClassRequest struct {
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type UpdateClassRequest struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}
