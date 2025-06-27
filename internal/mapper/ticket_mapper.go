package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Ticket domain to ReadTicketResponse model
func TicketToResponse(ticket *domain.Ticket) *model.ReadTicketResponse {
	return &model.ReadTicketResponse{
		ID: ticket.ID,
		Schedule: model.TicketSchedule{
			ID: ticket.Schedule.ID,
			Ship: model.TicketScheduleShip{
				ID:       ticket.Schedule.Ship.ID,
				ShipName: ticket.Schedule.Ship.ShipName,
			},
			DepartureHarbor: model.TicketScheduleHarbor{
				ID:         ticket.Schedule.DepartureHarbor.ID,
				HarborName: ticket.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.TicketScheduleHarbor{
				ID:         ticket.Schedule.ArrivalHarbor.ID,
				HarborName: ticket.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: ticket.Schedule.DepartureDatetime,
			ArrivalDatetime:   ticket.Schedule.ArrivalDatetime,
		},
		Class: model.TicketClassItem{
			ID:        ticket.Class.ID,
			ClassName: ticket.Class.ClassName,
			Type:      ticket.Class.Type,
		},
		TicketCode:      ticket.TicketCode, // Unique ticket code
		BookingID:       ticket.BookingID,
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
		Booking: &model.TicketBooking{
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
