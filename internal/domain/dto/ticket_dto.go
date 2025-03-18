package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type TicketHarborRes struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type TicketRouteRes struct {
	ID              uint            `json:"id"`
	DepartureHarbor TicketHarborRes `json:"departure_harbor"`
	ArrivalHarbor   TicketHarborRes `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type TicketShipRes struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
	Capacity uint   `json:"capacity"`
}

// ScheduleDTO represents a trip schedule.
type TicketScheduleRes struct {
	ID       uint           `json:"id"`
	DateTime time.Time      `json:"datetime"`
	Ship     TicketShipRes  `json:"ship"`
	Route    TicketRouteRes `json:"route"`
}

// ClassDTO represents ticket class information.
type TicketClassRes struct {
	ID        uint    `json:"id"`
	ClassName string  `json:"class_name"`
	Price     float64 `json:"price"`
}

// BookingDTO represents the person who booked the ticket.
type TicketBookingRes struct {
	ID          uint              `json:"id"`
	CusName     string            `json:"cus_name"`
	PersonID    uint              `json:"person_id"`
	PhoneNumber uint              `json:"phone_number"` // Changed to string to support leading zeros
	Email       string            `json:"email_address"`
	BirthDate   time.Time         `json:"birth_date"`
	Schedule    TicketScheduleRes `json:"schedule"`
}

// TicketDTO represents a ticket.
type TicketRes struct {
	ID            uint             `json:"id"`
	PassengerName string           `json:"passenger_name"`
	SeatNumber    string           `json:"seat_number"`
	Class         TicketClassRes   `json:"class"`
	Booking       TicketBookingRes `json:"booking"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func ToTicketDTO(ticket *entities.Ticket) TicketRes {
	var ticketResponse TicketRes
	copier.Copy(&ticketResponse, &ticket) // Automatically maps matching fields
	return ticketResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToTicketDTOs(tickets []*entities.Ticket) []TicketRes {
	var ticketResponses []TicketRes
	for _, ticket := range tickets {
		ticketResponses = append(ticketResponses, ToTicketDTO(ticket))
	}
	return ticketResponses
}
