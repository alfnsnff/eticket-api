package requests

import (
	"eticket-api/internal/domain"
	"time"
)

type CreateTicketRequest struct {
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

// Map Ticket domain to ReadTicketResponse model
func TicketToResponse(ticket *domain.Ticket) *TicketResponse {
	return &TicketResponse{
		ID: ticket.ID,
		Schedule: TicketSchedule{
			ID: ticket.Schedule.ID,
			Ship: TicketScheduleShip{
				ID:       ticket.Schedule.Ship.ID,
				ShipName: ticket.Schedule.Ship.ShipName,
			},
			DepartureHarbor: TicketScheduleHarbor{
				ID:         ticket.Schedule.DepartureHarbor.ID,
				HarborName: ticket.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: TicketScheduleHarbor{
				ID:         ticket.Schedule.ArrivalHarbor.ID,
				HarborName: ticket.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: ticket.Schedule.DepartureDatetime,
			ArrivalDatetime:   ticket.Schedule.ArrivalDatetime,
		},
		Class: TicketClass{
			ID:        ticket.Class.ID,
			ClassName: ticket.Class.ClassName,
			Type:      ticket.Class.Type,
		},
		TicketCode:      ticket.TicketCode, // Unique ticket code
		PassengerName:   ticket.PassengerName,
		PassengerAge:    ticket.PassengerAge,
		PassengerGender: ticket.PassengerGender,
		Address:         ticket.Address,
		IDType:          ticket.IDType,
		IDNumber:        ticket.IDNumber,
		SeatNumber:      ticket.SeatNumber,
		LicensePlate:    ticket.LicensePlate,
		IsCheckedIn:     ticket.IsCheckedIn,
		Type:            ticket.Type,
		Price:           ticket.Price,
		Booking: &TicketBooking{
			ID:             ticket.Booking.ID,
			OrderID:        ticket.Booking.OrderID, // Unique identifier for the booking
			CustomerName:   ticket.Booking.CustomerName,
			CustomerAge:    ticket.Booking.CustomerAge,
			CUstomerGender: ticket.Booking.CustomerGender,
			IDType:         ticket.Booking.IDType,
			IDNumber:       ticket.Booking.IDNumber,
			PhoneNumber:    ticket.Booking.PhoneNumber, // Changed to string to support leading zeros
			Email:          ticket.Booking.Email,
		},
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
	}

}

func TicketFromCreate(request *CreateTicketRequest) *domain.Ticket {
	return &domain.Ticket{
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		PassengerName:   request.PassengerName,
		PassengerAge:    request.PassengerAge,
		PassengerGender: request.PassengerGender,
		Address:         request.Address,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		SeatNumber:      request.SeatNumber,
		LicensePlate:    request.LicensePlate,
		Type:            request.Type,
		Price:           request.Price,
	}
}
func TicketFromUpdate(request *UpdateTicketRequest) *domain.Ticket {
	return &domain.Ticket{
		ID:              request.ID,
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		PassengerName:   request.PassengerName,
		PassengerAge:    request.PassengerAge,
		PassengerGender: request.PassengerGender,
		Address:         request.Address,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		SeatNumber:      request.SeatNumber,
		LicensePlate:    request.LicensePlate,
		Type:            request.Type,
		Price:           request.Price,
	}
}
