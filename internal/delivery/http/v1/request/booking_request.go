package request

import (
	"time"
)

type CreateBookingRequest struct {
	OrderID         string  `json:"order_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	IDType          string  `json:"id_type" validate:"required"`
	IDNumber        string  `json:"id_number" validate:"required"`
	CustomerName    string  `json:"customer_name" validate:"required"`
	CustomerAge     int     `json:"customer_age" validate:"required,min=0,max=120"`
	CustomerGender  string  `json:"customer_gender" validate:"required,oneof=male female other"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	ReferenceNumber *string `json:"reference_number"`
}

type UpdateBookingRequest struct {
	ID              uint    `json:"id,omitempty"`
	OrderID         string  `json:"order_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	CustomerName    string  `json:"customer_name" validate:"required"`
	CustomerAge     int     `json:"customer_age" validate:"required,min=0,max=120"`
	CustomerGender  string  `json:"customer_gender" validate:"required,oneof=male female other"`
	IDType          string  `json:"id_type" validate:"required"`
	IDNumber        string  `json:"id_number" validate:"required"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	ReferenceNumber *string `json:"reference_number"`
}

type BookingResponse struct {
	ID              uint            `json:"id"`
	OrderID         string          `json:"order_id"`
	Schedule        BookingSchedule `json:"schedule"`
	CustomerName    string          `json:"customer_name"`
	CustomerAge     int             `json:"customer_age"`
	CUstomerGender  string          `json:"customer_gender"`
	IDType          string          `json:"id_type"`
	IDNumber        string          `json:"id_number"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	ReferenceNumber *string         `json:"reference_number"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Tickets         []BookingTicket `json:"tickets"`
}

type BookingTicketClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}

type BookingTicket struct {
	ID              uint               `json:"id"`
	TicketCode      string             `json:"ticket_code"`
	Class           BookingTicketClass `json:"class"`
	Type            string             `json:"type" binding:"required,oneof=passenger vehicle"`
	PassengerName   string             `json:"passenger_name"`
	PassengerAge    int                `json:"passenger_age"`
	Address         string             `json:"address"`
	PassengerGender *string            `json:"passenger_gender"`
	IDType          *string            `json:"id_type"`
	IDNumber        *string            `json:"id_number"`
	SeatNumber      *string            `json:"seat_number"`
	LicensePlate    *string            `json:"license_plate"`
	Price           float64            `json:"price"`
}

type BookingHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type BookingScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type BookingSchedule struct {
	ID                uint                `json:"id"`
	Ship              BookingScheduleShip `json:"ship"`
	DepartureHarbor   BookingHarbor       `json:"departure_harbor"`
	ArrivalHarbor     BookingHarbor       `json:"arrival_harbor"`
	DepartureDatetime time.Time           `json:"departure_datetime"`
	ArrivalDatetime   time.Time           `json:"arrival_datetime"`
}
