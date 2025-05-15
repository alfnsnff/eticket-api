package model

import "time"

type ReadUserRoleResponse struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"role_id"` // e.g., "admin", "editor"
	RoleID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteUserRoleRequest struct {
	UserID string `json:"role_id"` // e.g., "admin", "editor"
	RoleID string `json:"user_id"`
}

// AllocationDTO represents a Allocation.
type UpdateUserRoleRequest struct {
	ID     uint   `json:"id"`
	UserID string `json:"role_id"` // e.g., "admin", "editor"
	RoleID string `json:"user_id"`
}
