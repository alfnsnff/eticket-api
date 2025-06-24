package ticket

import (
	"eticket-api/internal/common/utils"
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map Ticket domain to ReadTicketResponse model
func TicketToResponse(ticket *domain.Ticket) *model.ReadTicketResponse {
	return &model.ReadTicketResponse{
		ID:             ticket.ID,
		ClaimSessionID: ticket.ClaimSessionID,
		Schedule: model.TicketSchedule{
			ID: ticket.Schedule.ID,
			Ship: model.TicketScheduleShip{
				ID:       ticket.Schedule.Ship.ID,
				ShipName: ticket.Schedule.Ship.ShipName,
			},
			Route: model.TicketScheduleRoute{
				ID: ticket.Schedule.Route.ID,
				DepartureHarbor: model.TicketScheduleHarbor{
					ID:         ticket.Schedule.Route.DepartureHarbor.ID,
					HarborName: ticket.Schedule.Route.DepartureHarbor.HarborName,
				},
				ArrivalHarbor: model.TicketScheduleHarbor{
					ID:         ticket.Schedule.Route.ArrivalHarbor.ID,
					HarborName: ticket.Schedule.Route.ArrivalHarbor.HarborName,
				},
			},
			DepartureDatetime: utils.SafeTimeDeref(ticket.Schedule.DepartureDatetime),
			ArrivalDatetime:   utils.SafeTimeDeref(ticket.Schedule.ArrivalDatetime),
		},
		Class: model.TicketClassItem{
			ID:        ticket.Class.ID,
			ClassName: ticket.Class.ClassName,
			Type:      ticket.Class.Type,
		},
		BookingID:     ticket.BookingID,
		PassengerName: ticket.PassengerName,
		PassengerAge:  ticket.PassengerAge,
		Address:       ticket.Address,
		IDType:        ticket.IDType,
		IDNumber:      ticket.IDNumber,
		SeatNumber:    ticket.SeatNumber,
		LicensePlate:  ticket.LicensePlate,
		IsCheckedIn:   ticket.IsCheckedIn,
		Type:          ticket.Type,
		Price:         ticket.Price,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}

}
