package request

import (
	"time"
)

type WriteTicketRequest struct {
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	BookingID       *uint   `json:"booking_id"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	Address         string  `json:"address"`
	PassengerGender *string `json:"passenger_gender"`
	IDType          *string `json:"id_type"`
	IDNumber        *string `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Type            string  `json:"type" validate:"required"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	IsCheckedIn     bool    `json:"is_checked_in"`
}

type UpdateTicketRequest struct {
	ID              uint    `json:"id" validate:"required"`
	BookingID       *uint   `json:"booking_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	ClassID         uint    `json:"class_id" validate:"required"`
	PassengerName   string  `json:"passenger_name"`
	PassengerAge    int     `json:"passenger_age"`
	Address         string  `json:"address"`
	PassengerGender *string `json:"passenger_gender"`
	IDType          *string `json:"id_type"`
	IDNumber        *string `json:"id_number"`
	SeatNumber      *string `json:"seat_number"`
	LicensePlate    *string `json:"license_plate"`
	Type            string  `json:"type" validate:"required"`
	Price           float64 `json:"price" validate:"required,gte=0"`
	IsCheckedIn     bool    `json:"is_checked_in"`
}

type TicketResponse struct {
	ID              uint           `json:"id"`
	Schedule        TicketSchedule `json:"schedule"`
	Class           TicketClass    `json:"class"`
	TicketCode      string         `json:"ticket_code"`
	BookingID       *uint          `json:"booking_id"`
	Booking         *TicketBooking `json:"booking,omitempty"`
	PassengerName   string         `json:"passenger_name"`
	PassengerAge    int            `json:"passenger_age"`
	Address         string         `json:"address"`
	PassengerGender *string        `json:"passenger"`
	IDType          *string        `json:"id_type"`
	IDNumber        *string        `json:"id_number"`
	SeatNumber      *string        `json:"seat_number"`
	LicensePlate    *string        `json:"license_plate"`
	Type            string         `json:"type" binding:"required,oneof=passenger vehicle"`
	Price           float64        `json:"price"`
	IsCheckedIn     bool           `json:"is_checked_in"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type TicketScheduleHarbor struct {
	ID         uint   `json:"id"`
	HarborName string `json:"harbor_name"`
}

type TicketScheduleShip struct {
	ID       uint   `json:"id"`
	ShipName string `json:"ship_name"`
}

type TicketSchedule struct {
	ID                uint                 `json:"id"`
	Ship              TicketScheduleShip   `json:"ship"`
	DepartureHarbor   TicketScheduleHarbor `json:"departure_harbor"`
	ArrivalHarbor     TicketScheduleHarbor `json:"arrival_harbor"`
	DepartureDatetime time.Time            `json:"departure_datetime"`
	ArrivalDatetime   time.Time            `json:"arrival_datetime"`
}

type TicketBooking struct {
	ID             uint   `json:"id"`
	OrderID        string `json:"order_id"`
	CustomerName   string `json:"customer_name"`
	CustomerAge    int    `json:"customer_age"`
	CUstomerGender string `json:"customer_gender"`
	IDType         string `json:"id_type"`
	IDNumber       string `json:"id_number"`
	PhoneNumber    string `json:"phone_number"`
	Email          string `json:"email"`
}

type TicketClass struct {
	ID        uint   `json:"id"`
	ClassName string `json:"class_name"`
	Type      string `json:"type"`
}
