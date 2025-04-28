package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type HarborRead struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HarborDTO represents a harbor.
type HarborCreate struct {
	Name string `json:"name"`
}

func ToHarborDTO(harbor *entities.Harbor) HarborRead {
	var harborResponse HarborRead
	copier.Copy(&harborResponse, &harbor) // Automatically maps matching fields
	return harborResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToHarborDTOs(harbors []*entities.Harbor) []HarborRead {
	var harborResponses []HarborRead
	for _, harbor := range harbors {
		harborResponses = append(harborResponses, ToHarborDTO(harbor))
	}
	return harborResponses
}

func ToHarborEntity(harborCreate *HarborCreate) entities.Harbor {
	var harbor entities.Harbor
	copier.Copy(&harbor, &harborCreate)
	return harbor
}
