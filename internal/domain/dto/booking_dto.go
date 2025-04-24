package dto

import (
	"eticket-api/internal/domain/entities"
	"time"

	"github.com/jinzhu/copier"
)

// HarborDTO represents a harbor.
type BookingHarbor struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// RouteDTO represents a travel route.
type BookingRoute struct {
	ID              uint          `json:"id"`
	DepartureHarbor BookingHarbor `json:"departure_harbor"`
	ArrivalHarbor   BookingHarbor `json:"arrival_harbor"`
}

// ShipDTO represents a ship.
type BookingShip struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ScheduleDTO represents a trip schedule.
type BookingSchedule struct {
	ID       uint         `json:"id"`
	DateTime time.Time    `json:"datetime"`
	Ship     BookingShip  `json:"ship"`
	Route    BookingRoute `json:"route"`
}

// ClassDTO represents ticket class information.
type BookingTicketClass struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookingShipClass struct {
	ID    uint        `json:"id"`
	Class TicketClass `json:"class"`
}

type BookingTicketPrice struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Price     float32          `json:"price"`
	ShipClass BookingShipClass `json:"ship_class"`
}

type BookingTicketCreate struct {
	ID            uint   `json:"id"`
	PriceID       uint   `json:"price_id"`
	ScheduleID    uint   `json:"schedule_id"`
	PassengerName string `json:"passenger_name"`
	SeatNumber    string `json:"seat_number"`
}

type BookingTicketRead struct {
	PassengerName string      `json:"passenger_name"`
	SeatNumber    string      `json:"seat_number"`
	Price         TicketPrice `json:"price"`
}

// BookingDTO represents the person who booked the ticket.
type BookingRead struct {
	ID          uint                `json:"id"`
	CusName     string              `json:"cus_name"`
	PersonID    uint                `json:"person_id"`
	PhoneNumber string              `json:"phone_number"` // Changed to string to support leading zeros
	Email       string              `json:"email_address"`
	BirthDate   time.Time           `json:"birth_date"`
	Schedule    BookingSchedule     `json:"schedule"`
	Tickets     []BookingTicketRead `json:"tickets"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type BookingCreate struct {
	ScheduleID  uint                  `json:"schedule_id"` // Foreign key
	CusName     string                `json:"cus_name"`
	PersonID    uint                  `json:"person_id"`
	PhoneNumber string                `json:"phone_number"` // Changed to string to support leading zeros
	Email       string                `json:"email_address"`
	BirthDate   time.Time             `json:"birth_date"`
	Schedule    BookingSchedule       `json:"schedule"`
	Tickets     []BookingTicketCreate `json:"tickets"`
}

func ToBookingDTO(booking *entities.Booking) BookingRead {
	var bookingRead BookingRead
	copier.Copy(&bookingRead, &booking) // Automatically maps matching fields
	return bookingRead
}

// Convert a slice of Ticket entities to DTO slice
func ToBookingDTOs(bookings []*entities.Booking) []BookingRead {
	var bookingRead []BookingRead
	for _, booking := range bookings {
		bookingRead = append(bookingRead, ToBookingDTO(booking))
	}
	return bookingRead
}

// Convert BookingReq DTO to Booking entity
func ToBookingEntity(bookingCreate *BookingCreate) (entities.Booking, []entities.Ticket) {
	var booking entities.Booking
	var tickets []entities.Ticket

	// Automatically copy matching fields
	copier.Copy(&booking, &bookingCreate)

	// Convert tickets manually (copier doesn't automatically handle slices of nested structs)
	for _, t := range bookingCreate.Tickets {
		var ticket entities.Ticket
		copier.Copy(&ticket, &t)
		tickets = append(tickets, ticket)
	}

	return booking, tickets
}
