package model

import "time"

type ReadRoleResponse struct {
	ID          uint      `json:"id"`
	RoleName    string    `json:"role_name"` // e.g., "admin", "editor"
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteRoleRequest struct {
	RoleName    string `json:"role_name"` // e.g., "admin", "editor"
	Description string `json:"description"`
}

// AllocationDTO represents a Allocation.
type UpdateRoleRequest struct {
	ID          uint   `json:"id"`
	RoleName    string `json:"role_name"` // e.g., "admin", "editor"
	Description string `json:"description"`
}
