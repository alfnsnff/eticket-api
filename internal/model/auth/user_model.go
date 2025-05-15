package model

import "time"

type ReadUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"` // e.g., "admin", "editor"
	Email     string    `json:"email"`
	Fullname  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AllocationDTO represents a Allocation.
type WriteUserRequest struct {
	Username string `json:"username"` // e.g., "admin", "editor"
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
}

// AllocationDTO represents a Allocation.
type UpdateuserRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username"` // e.g., "admin", "editor"
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"full_name"`
}
