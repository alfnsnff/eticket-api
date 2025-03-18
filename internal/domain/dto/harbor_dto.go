package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type HarborRes struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	HarborName string    `gorm:"not null" json:"harbor_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToHarborDTO(harbor *entities.Harbor) HarborRes {
	var harborResponse HarborRes
	copier.Copy(&harborResponse, &harbor) // Automatically maps matching fields
	return harborResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToHarborDTOs(harbors []*entities.Harbor) []HarborRes {
	var harborResponses []HarborRes
	for _, harbor := range harbors {
		harborResponses = append(harborResponses, ToHarborDTO(harbor))
	}
	return harborResponses
}
