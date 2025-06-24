package model

import (
	"time"
)

// HarborDTO represents a harbor.
type ReadHarborResponse struct {
	ID            uint      `json:"id"`
	HarborName    string    `json:"harbor_name"`
	Status        string    `json:"status"`
	HarborAlias   *string   `json:"harbor_alias,omitempty"`
	YearOperation string    `json:"year_operation"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type WriteHarborRequest struct {
	HarborName    string  `json:"harbor_name" validate:"required"`
	Status        string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
	YearOperation string  `json:"year_operation" validate:"required,len=4,numeric"`
	HarborAlias   *string `json:"harbor_alias,omitempty"` // Optional
}

type UpdateHarborRequest struct {
	ID            uint    `json:"id" validate:"required"`
	HarborName    string  `json:"harbor_name" validate:"required"`
	Status        string  `json:"status" validate:"required,oneof=ACTIVE INACTIVE"`
	YearOperation string  `json:"year_operation" validate:"required,len=4,numeric"`
	HarborAlias   *string `json:"harbor_alias,omitempty"` // Optional
}
