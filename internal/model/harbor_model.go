package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ReadHarborResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Year      string    `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HarborDTO represents a harbor.
type WriteHarborRequest struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Year   string `json:"year"`
}

// HarborDTO represents a harbor.
type UpdateHarborRequest struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Year   string `json:"year"`
}
