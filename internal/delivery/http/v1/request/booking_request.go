package requests

import (
	"eticket-api/internal/domain"
	"time"
)

type CreateBookingRequest struct {
	OrderID         string  `json:"order_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	IDType          string  `json:"id_type" validate:"required"`
	IDNumber        string  `json:"id_number" validate:"required"`
	CustomerName    string  `json:"customer_name" validate:"required"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	Status          string  `json:"status" validate:"required,oneof=paid unpaid expired refunded"`
	ReferenceNumber *string `json:"reference_number"`
}

type UpdateBookingRequest struct {
	ID              uint    `json:"id" validate:"required"`
	OrderID         string  `json:"order_id"`
	ScheduleID      uint    `json:"schedule_id" validate:"required"`
	CustomerName    string  `json:"customer_name" validate:"required"`
	IDType          string  `json:"id_type" validate:"required"`
	IDNumber        string  `json:"id_number" validate:"required"`
	PhoneNumber     string  `json:"phone_number" validate:"required"`
	Email           string  `json:"email" validate:"required,email"`
	Status          string  `json:"status" validate:"required,oneof=paid unpaid expired refunded"`
	ReferenceNumber *string `json:"reference_number"`
}

type RefundBookingRequest struct {
	OrderID  string `json:"order_id"`
	IDType   string `json:"id_type" validate:"required"`
	IDNumber string `json:"id_number" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type BookingResponse struct {
	ID              uint            `json:"id"`
	OrderID         string          `json:"order_id"`
	Schedule        BookingSchedule `json:"schedule"`
	CustomerName    string          `json:"customer_name"`
	IDType          string          `json:"id_type"`
	IDNumber        string          `json:"id_number"`
	PhoneNumber     string          `json:"phone_number"`
	Email           string          `json:"email"`
	Status          string          `json:"status"`
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

// Map Booking domain to ReadBookingResponse model
func BookingToResponse(booking *domain.Booking) *BookingResponse {
	// Map tickets
	tickets := make([]BookingTicket, len(booking.Tickets))
	for i, ticket := range booking.Tickets {
		tickets[i] = BookingTicket{
			ID:   ticket.ID,
			Type: ticket.Type,
			Class: BookingTicketClass{
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
			Price:           ticket.Price,
		}
	}

	return &BookingResponse{
		ID:      booking.ID,
		OrderID: booking.OrderID,
		Schedule: BookingSchedule{
			ID: booking.Schedule.ID,
			Ship: BookingScheduleShip{
				ID:       booking.Schedule.Ship.ID,
				ShipName: booking.Schedule.Ship.ShipName,
			},
			DepartureHarbor: BookingHarbor{
				ID:         booking.Schedule.DepartureHarbor.ID,
				HarborName: booking.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: BookingHarbor{
				ID:         booking.Schedule.ArrivalHarbor.ID,
				HarborName: booking.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: booking.Schedule.DepartureDatetime,
			ArrivalDatetime:   booking.Schedule.ArrivalDatetime,
		},
		CustomerName:    booking.CustomerName,
		IDType:          booking.IDType,
		IDNumber:        booking.IDNumber,
		PhoneNumber:     booking.PhoneNumber,
		Email:           booking.Email,
		Status:          booking.Status,
		ReferenceNumber: booking.ReferenceNumber,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		Tickets:         tickets,
	}
}

func BookingFromCreate(request *CreateBookingRequest) *domain.Booking {
	return &domain.Booking{
		OrderID:         request.OrderID,
		ScheduleID:      request.ScheduleID,
		CustomerName:    request.CustomerName,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		PhoneNumber:     request.PhoneNumber,
		Email:           request.Email,
		Status:          request.Status,
		ReferenceNumber: request.ReferenceNumber,
	}
}

func BookingFromUpdate(request *UpdateBookingRequest) *domain.Booking {
	return &domain.Booking{
		ID:              request.ID,
		OrderID:         request.OrderID,
		ScheduleID:      request.ScheduleID,
		CustomerName:    request.CustomerName,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		PhoneNumber:     request.PhoneNumber,
		Email:           request.Email,
		Status:          request.Status,
		ReferenceNumber: request.ReferenceNumber,
	}
}
