package ticket

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
	"time"
)

// Map Ticket entity to ReadTicketResponse model
func ToReadTicketResponse(ticket *entity.Ticket) *model.ReadTicketResponse {
	if ticket == nil {
		return nil
	}

	return &model.ReadTicketResponse{
		ID:             ticket.ID,
		ClaimSessionID: *ticket.ClaimSessionID,
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
			DepartureDatetime: *ticket.Schedule.DepartureDatetime,
			ArrivalDatetime:   *ticket.Schedule.ArrivalDatetime,
		},
		Class: model.TicketClassItem{
			ID:        ticket.Class.ID,
			ClassName: ticket.Class.ClassName,
			Type:      ticket.Class.Type,
		},
		Status:        "active", // Example status
		BookingID:     *ticket.BookingID,
		Type:          ticket.Type,
		PassengerName: *ticket.PassengerName,
		PassengerAge:  *ticket.PassengerAge,
		Address:       *ticket.Address,
		IDType:        *ticket.IDType,
		IDNumber:      *ticket.IDNumber,
		SeatNumber:    ticket.SeatNumber,
		LicensePlate:  ticket.LicensePlate,
		Price:         ticket.Price,
		ExpiresAt:     ticket.CreatedAt.Add(24 * time.Hour), // Example expiration logic
		ClaimedAt:     ticket.UpdatedAt,                     // Example claimed logic
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
	}
}

// Map a slice of Ticket entities to ReadTicketResponse models
func ToReadTicketResponses(tickets []*entity.Ticket) []*model.ReadTicketResponse {
	responses := make([]*model.ReadTicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = ToReadTicketResponse(ticket)
	}
	return responses
}

// Map WriteTicketRequest model to Ticket entity
func FromWriteTicketRequest(request *model.WriteTicketRequest) *entity.Ticket {
	return &entity.Ticket{
		ScheduleID:      request.ScheduleID,
		ClassID:         request.ClassID,
		BookingID:       request.BookingID,
		ClaimSessionID:  request.ClaimSessionID,
		Type:            request.Type,
		Price:           request.Price,
		PassengerName:   request.PassengerName,
		PassengerAge:    request.PassengerAge,
		PassengerGender: request.PassengerGender,
		Address:         request.Address,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		SeatNumber:      request.SeatNumber,
		LicensePlate:    request.LicensePlate,
	}
}

// Map UpdateTicketRequest model to Ticket entity
func FromUpdateTicketRequest(request *model.UpdateTicketRequest, ticket *entity.Ticket) {
	ticket.ScheduleID = request.ScheduleID
	ticket.ClassID = request.ClassID
	ticket.BookingID = &request.BookingID
	ticket.ClaimSessionID = &request.ClaimSessionID
	ticket.Type = request.Type
	ticket.Price = request.Price
	ticket.PassengerName = &request.PassengerName
	ticket.PassengerAge = &request.PassengerAge
	ticket.PassengerGender = &request.PassengerGender
	ticket.Address = &request.Address
	ticket.IDType = &request.IDType
	ticket.IDNumber = &request.IDNumber
	ticket.SeatNumber = request.SeatNumber
	ticket.LicensePlate = request.LicensePlate
}
