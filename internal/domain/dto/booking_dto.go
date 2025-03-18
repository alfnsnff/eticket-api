package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type BookingHarborRes struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

// RouteDTO represents a travel route.
type BookingRouteRes struct {
	ID              uint             `json:"id"`
	DepartureHarbor BookingHarborRes `json:"departure_harbor"`
	ArrivalHarbor   BookingHarborRes `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type BookingShipRes struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
	Capacity uint   `json:"capacity"`
}

// ScheduleDTO represents a trip schedule.
type BookingScheduleRes struct {
	ID       uint            `json:"id"`
	DateTime time.Time       `json:"datetime"`
	Ship     BookingShipRes  `json:"ship"`
	Route    BookingRouteRes `json:"route"`
}

// BookingDTO represents the person who booked the ticket.
type BookingRes struct {
	ID          uint               `json:"id"`
	CusName     string             `json:"cus_name"`
	PersonID    uint               `json:"person_id"`
	PhoneNumber uint               `json:"phone_number"` // Changed to string to support leading zeros
	Email       string             `json:"email_address"`
	BirthDate   time.Time          `json:"birth_date"`
	Schedule    BookingScheduleRes `json:"schedule"`
}

func ToBookingDTO(booking *entities.Booking) BookingRes {
	var bookingResponse BookingRes
	copier.Copy(&bookingResponse, &booking) // Automatically maps matching fields
	return bookingResponse
}

// Convert a slice of Ticket entities to DTO slice
func ToBookingDTOs(bookings []*entities.Booking) []BookingRes {
	var bookingResponses []BookingRes
	for _, booking := range bookings {
		bookingResponses = append(bookingResponses, ToBookingDTO(booking))
	}
	return bookingResponses
}
