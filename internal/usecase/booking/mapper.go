package booking

import (
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
)

// Map Booking domain to ReadBookingResponse model
func BookingToResponse(booking *domain.Booking) *model.ReadBookingResponse {

	// Print tickets to terminal
	for _, ticket := range booking.Tickets {
		fmt.Printf("Ticket: %+v\n", ticket)
	}

	// Map tickets
	tickets := make([]model.BookingTicket, len(booking.Tickets))
	for i, ticket := range booking.Tickets {
		tickets[i] = model.BookingTicket{
			ID:            ticket.ID,
			Type:          ticket.Type,
			PassengerName: ticket.PassengerName,
			PassengerAge:  ticket.PassengerAge,
			Address:       ticket.Address,
			IDType:        ticket.IDType,
			IDNumber:      ticket.IDNumber,
			SeatNumber:    ticket.SeatNumber,
			LicensePlate:  ticket.LicensePlate,
			Price:         ticket.Price,
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
			Route: model.BookingScheduleRoute{
				ID: booking.Schedule.Route.ID,
				DepartureHarbor: model.BookingScheduleHarbor{
					ID:         booking.Schedule.Route.DepartureHarbor.ID,
					HarborName: booking.Schedule.Route.DepartureHarbor.HarborName,
				},
				ArrivalHarbor: model.BookingScheduleHarbor{
					ID:         booking.Schedule.Route.ArrivalHarbor.ID,
					HarborName: booking.Schedule.Route.ArrivalHarbor.HarborName,
				},
			},
			DepartureDatetime: utils.SafeTimeDeref(booking.Schedule.DepartureDatetime),
			ArrivalDatetime:   utils.SafeTimeDeref(booking.Schedule.ArrivalDatetime),
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
