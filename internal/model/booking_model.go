package model

import (
	"time"
)

// // HarborDTO represents a harbor.
// type BookingHarbor struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // RouteDTO represents a travel route.
// type BookingRoute struct {
// 	ID              uint          `json:"id"`
// 	DepartureHarbor BookingHarbor `json:"departure_harbor"`
// 	ArrivalHarbor   BookingHarbor `json:"arrival_harbor"`
// }

// // ShipDTO represents a ship.
// type BookingShip struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// // ScheduleDTO represents a trip schedule.
// type BookingSchedule struct {
// 	ID       uint         `json:"id"`
// 	DateTime time.Time    `json:"datetime"`
// 	Ship     BookingShip  `json:"ship"`
// 	Route    BookingRoute `json:"route"`
// }

// // ClassDTO represents ticket class information.
// type BookingTicketClass struct {
// 	ID   uint   `json:"id"`
// 	Name string `json:"name"`
// }

// type BookingManifest struct {
// 	ID    uint        `json:"id"`
// 	Class TicketClass `json:"class"`
// }

// type BookingTicketFare struct {
// 	ID       uint            `json:"id"`
// 	Price    float32         `json:"price"`
// 	Manifest BookingManifest `json:"manifest"`
// }

// type WriteBookingTicketRequest struct {
// 	FareID        uint   `json:"fare_id"`
// 	ScheduleID    uint   `json:"schedule_id"`
// 	PassengerName string `json:"passenger_name"`
// 	SeatNumber    string `json:"seat_number"`
// }

// type ReadBookingTicketResponse struct {
// 	PassengerName string            `json:"passenger_name"`
// 	SeatNumber    string            `json:"seat_number"`
// 	Fare          BookingTicketFare `json:"fare"`
// }

// BookingDTO represents the person who booked the ticket.
type ReadBookingResponse struct {
	ID           uint      `json:"id"`
	CustomerName string    `json:"customer_name"`
	IDType       uint      `json:"id_type"`
	IDNumber     uint      `json:"id_number"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`

	BookingTimestamp time.Time `gorm:"not null" json:"booking_timestamp"`       // Timestamp when the booking was confirmed
	TotalAmount      float32   `gorm:"not null" json:"total_amount"`            // Total price of all tickets in this booking
	Status           string    `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WriteBookingRequest struct {
	ScheduleID   uint      `json:"schedule_id"` // Foreign key
	CustomerName string    `json:"customer_name"`
	PersonID     uint      `json:"person_id"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`
}

type UpdateBookingRequest struct {
	ID           uint      `json:"id"`
	ScheduleID   uint      `json:"schedule_id"` // Foreign key
	CustomerName string    `json:"customer_name"`
	PersonID     uint      `json:"person_id"`
	PhoneNumber  string    `json:"phone_number"` // Changed to string to support leading zeros
	Email        string    `json:"email_address"`
	BirthDate    time.Time `json:"birth_date"`
}

// type Booking struct {
// 	ID           uint      `gorm:"primaryKey" json:"id"`
// 	ScheduleID   uint      `gorm:"not null;index;" json:"schedule_id"` // Foreign key
// 	PersonID     uint      `gorm:"not null" json:"person_id"`
// 	IDType       string    `gorm:"type:varchar(10);not null" json:"id_type"`   // Changed to string to support leading zeros
// 	IDNumber     string    `gorm:"type:varchar(10);not null" json:"id_number"` // Changed to string to support leading zeros
// 	CustomerName string    `gorm:"not null" json:"customer_name"`
// 	PhoneNumber  string    `gorm:"type:varchar(15);not null" json:"phone_number"` // Changed to string to support leading zeros
// 	Email        string    `gorm:"not null" json:"email"`
// 	BirthDate    time.Time `gorm:"not null" json:"birth_date"`

// 	BookingTimestamp time.Time `gorm:"not null" json:"booking_timestamp"`       // Timestamp when the booking was confirmed
// 	TotalAmount      float32   `gorm:"not null" json:"total_amount"`            // Total price of all tickets in this booking
// 	Status           string    `gorm:"type:varchar(20);not null" json:"status"` // e.g., 'completed', 'cancelled', 'refunded'

// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// }
