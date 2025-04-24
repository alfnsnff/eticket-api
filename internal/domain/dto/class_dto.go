package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type ClassHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type ClassRoute struct {
	ID              uint        `json:"id"`
	DepartureHarbor ClassHarbor `json:"departure_harbor"`
	ArrivalHarbor   ClassHarbor `json:"arrival_harbor"`
}

// ClassDTO represents ticket class information.
type ClassRead struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ClassCreate struct {
	Name string `json:"name"`
}

func ToClassDTO(class *entities.Class) ClassRead {
	var classRead ClassRead
	copier.Copy(&classRead, &class) // Automatically maps matching fields
	return classRead
}

// Convert a slice of Ticket entities to DTO slice
func ToClassDTOs(classes []*entities.Class) []ClassRead {
	var classReads []ClassRead
	for _, class := range classes {
		classReads = append(classReads, ToClassDTO(class))
	}
	return classReads
}

func ToClassEntity(classCreate *ClassCreate) entities.Class {
	var class entities.Class
	copier.Copy(&class, &classCreate)
	return class
}
