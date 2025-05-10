package model

import (
	"time"
)

// ClassDTO represents ticket class information.
type ReadClassResponse struct {
	ID        uint      `json:"id"`
	ClassName string    `json:"class_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteClassRequest struct {
	ClassName string `json:"class_name"`
}

type UpdateClassRequest struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
}
