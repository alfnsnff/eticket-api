package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ReadHarborResponse struct {
	ID            uint      `json:"id"`
	HarborName    string    `json:"harbor_name"`
	Status        string    `json:"status"`
	YearOperation string    `json:"year_operation"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// HarborDTO represents a harbor.
type WriteHarborRequest struct {
	HarborName    string `json:"harbor_name"`
	Status        string `json:"status"`
	YearOperation string `json:"year_operation"`
}

// HarborDTO represents a harbor.
type UpdateHarborRequest struct {
	ID            uint   `json:"id"`
	HarborName    string `json:"harbor_name"`
	Status        string `json:"status"`
	YearOperation string `json:"year_operation"`
}
