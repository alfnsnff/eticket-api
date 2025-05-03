package model

import (
	"time"
)

// HarborDTO represents a harbor.
// type ClassHarbor struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // RouteDTO represents a travel route.
// type ClassRoute struct {
// 	ID              uint        `json:"id"`
// 	DepartureHarbor ClassHarbor `json:"departure_harbor"`
// 	ArrivalHarbor   ClassHarbor `json:"arrival_harbor"`
// }

// ClassDTO represents ticket class information.
type ReadClassResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteClassRequest struct {
	Name string `json:"name"`
}

type UpdateClassRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
