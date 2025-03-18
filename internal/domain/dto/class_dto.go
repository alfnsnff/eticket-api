package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type ClassHarborRes struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type ClassRouteRes struct {
	ID              uint           `json:"id"`
	DepartureHarbor ClassHarborRes `json:"departure_harbor"`
	ArrivalHarbor   ClassHarborRes `json:"arrival_harbor"`
}

// ClassDTO represents ticket class information.
type ClassRes struct {
	ID        uint          `json:"id"`
	ClassName string        `json:"class_name"`
	Price     float64       `json:"price"`
	Route     ClassRouteRes `json:"route"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func ToClassDTO(class *entities.Class) ClassRes {
	var classResponse ClassRes
	copier.Copy(&classResponse, &class) // Automatically maps matching fields
	return classResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToClassDTOs(classes []*entities.Class) []ClassRes {
	var classResponses []ClassRes
	for _, class := range classes {
		classResponses = append(classResponses, ToClassDTO(class))
	}
	return classResponses
}
