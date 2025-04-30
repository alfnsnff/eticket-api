package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ReadHarborResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HarborDTO represents a harbor.
type WriteHarborRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
