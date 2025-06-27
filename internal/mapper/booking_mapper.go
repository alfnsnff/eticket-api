package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Booking domain to ReadBookingResponse model
func BookingToResponse(booking *domain.Booking) *model.ReadBookingResponse {
	// Map tickets
	tickets := make([]model.BookingTicket, len(booking.Tickets))
	for i, ticket := range booking.Tickets {
		tickets[i] = model.BookingTicket{
			ID:   ticket.ID,
			Type: ticket.Type,
			Class: model.BookingTicketClass{
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

	return &model.ReadBookingResponse{
		ID:      booking.ID,
		OrderID: booking.OrderID,
		Schedule: model.BookingSchedule{
			ID: booking.Schedule.ID,
			Ship: model.BookingScheduleShip{
				ID:       booking.Schedule.Ship.ID,
				ShipName: booking.Schedule.Ship.ShipName,
			},
			DepartureHarbor: model.BookingHarbor{
				ID:         booking.Schedule.DepartureHarbor.ID,
				HarborName: booking.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.BookingHarbor{
				ID:         booking.Schedule.ArrivalHarbor.ID,
				HarborName: booking.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: booking.Schedule.DepartureDatetime,
			ArrivalDatetime:   booking.Schedule.ArrivalDatetime,
		},
		CustomerName:    booking.CustomerName,
		CustomerAge:     booking.CustomerAge,
		CUstomerGender:  booking.CustomerGender,
		IDType:          booking.IDType,
		IDNumber:        booking.IDNumber,
		PhoneNumber:     booking.PhoneNumber,
		Email:           booking.Email,
		ReferenceNumber: booking.ReferenceNumber,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		Tickets:         tickets,
	}
}
