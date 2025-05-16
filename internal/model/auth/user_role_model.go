package model

import "time"

type ReadUserRoleResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"role_id"` // e.g., "admin", "editor"
	RoleID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteUserRoleRequest struct {
	UserID uint `json:"role_id"` // e.g., "admin", "editor"
	RoleID uint `json:"user_id"`
}

// AllocationDTO represents a Allocation.
type UpdateUserRoleRequest struct {
	ID     uint `json:"id"`
	UserID uint `json:"role_id"` // e.g., "admin", "editor"
	RoleID uint `json:"user_id"`
}
