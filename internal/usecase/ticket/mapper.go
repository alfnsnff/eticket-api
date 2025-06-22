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

	response := &model.ReadTicketResponse{
		ID:   ticket.ID,
		Type: ticket.Type,
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
		},
		Class: model.TicketClassItem{
			ID:        ticket.Class.ID,
			ClassName: ticket.Class.ClassName,
			Type:      ticket.Class.Type,
		},
		Price:     ticket.Price,
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
	}

	// Handle nullable fields safely
	if ticket.ClaimSessionID != nil {
		response.ClaimSessionID = *ticket.ClaimSessionID
	}

	if ticket.BookingID != nil {
		response.BookingID = *ticket.BookingID
	}

	if ticket.PassengerName != nil {
		response.PassengerName = *ticket.PassengerName
	}

	if ticket.PassengerAge != nil {
		response.PassengerAge = *ticket.PassengerAge
	}

	if ticket.Address != nil {
		response.Address = *ticket.Address
	}

	if ticket.IDType != nil {
		response.IDType = *ticket.IDType
	}

	if ticket.IDNumber != nil {
		response.IDNumber = *ticket.IDNumber
	}

	if ticket.SeatNumber != nil {
		response.SeatNumber = ticket.SeatNumber
	}

	if ticket.LicensePlate != nil {
		response.LicensePlate = ticket.LicensePlate
	}

	// Handle schedule datetime fields
	if ticket.Schedule.DepartureDatetime != nil {
		response.Schedule.DepartureDatetime = *ticket.Schedule.DepartureDatetime
	}

	if ticket.Schedule.ArrivalDatetime != nil {
		response.Schedule.ArrivalDatetime = *ticket.Schedule.ArrivalDatetime
	}

	// Set expiration and claimed time logic
	response.ExpiresAt = ticket.CreatedAt.Add(24 * time.Hour)
	response.ClaimedAt = ticket.UpdatedAt

	return response
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
