package claim_session

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map ClaimSession entity to ReadClaimSessionResponse model
func ToReadClaimSessionResponse(session *entity.ClaimSession, tickets []*entity.Ticket) *model.ReadClaimSessionResponse {
	if session == nil {
		return nil
	}

	var ticketPrices []model.ClaimSessionTicketPricesResponse
	var ticketDetails []model.ClaimSessionTicketDetailResponse
	var total float32

	// Build ticket prices and details
	if len(tickets) > 0 {
		ticketPrices, total = BuildPriceBreakdown(tickets)
		ticketDetails = BuildTicketBreakdown(tickets)
	}

	// Map schedule
	var scheduleModel model.ClaimSessionSchedule
	if session.Schedule.ID != 0 {
		scheduleModel = ToScheduleModel(&session.Schedule)
	}

	return &model.ReadClaimSessionResponse{
		ID:          session.ID,
		SessionID:   session.SessionID,
		ScheduleID:  session.ScheduleID,
		Schedule:    scheduleModel,
		ExpiresAt:   session.ExpiresAt,
		Prices:      ticketPrices,
		Tickets:     ticketDetails,
		TotalAmount: total,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// Map Schedule entity to ClaimSessionSchedule model
func ToScheduleModel(schedule *entity.Schedule) model.ClaimSessionSchedule {
	return model.ClaimSessionSchedule{
		ID: schedule.ID,
		Ship: model.ClaimSessionScheduleShip{
			ID:       schedule.Ship.ID,
			ShipName: schedule.Ship.ShipName,
		},
		Route: model.ClaimSessionScheduleRoute{
			ID: schedule.Route.ID,
			DepartureHarbor: model.ClaimSessionScheduleHarbor{
				ID:         schedule.Route.DepartureHarbor.ID,
				HarborName: schedule.Route.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.ClaimSessionScheduleHarbor{
				ID:         schedule.Route.ArrivalHarbor.ID,
				HarborName: schedule.Route.ArrivalHarbor.HarborName,
			},
		},
		DepartureDatetime: *schedule.DepartureDatetime,
		ArrivalDatetime:   *schedule.ArrivalDatetime,
	}
}

// Build ticket breakdown
func BuildTicketBreakdown(tickets []*entity.Ticket) []model.ClaimSessionTicketDetailResponse {
	result := make([]model.ClaimSessionTicketDetailResponse, len(tickets))
	for i, ticket := range tickets {
		result[i] = model.ClaimSessionTicketDetailResponse{
			TicketID: ticket.ID,
			Class: model.ClaimSessionTicketClassItem{
				ID:        ticket.Class.ID,
				ClassName: ticket.Class.ClassName,
				Type:      ticket.Class.Type,
			},
			Price: ticket.Price,
			Type:  ticket.Type,
		}
	}
	return result
}

// Build price breakdown
func BuildPriceBreakdown(tickets []*entity.Ticket) ([]model.ClaimSessionTicketPricesResponse, float32) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float32

	for _, ticket := range tickets {
		classID := ticket.ClassID
		price := ticket.Price

		if _, exists := ticketSummary[classID]; !exists {
			ticketSummary[classID] = &model.ClaimSessionTicketPricesResponse{
				Class: model.ClaimSessionTicketClassItem{
					ID:        ticket.Class.ID,
					ClassName: ticket.Class.ClassName,
					Type:      ticket.Class.Type,
				},
				Price:    price,
				Quantity: 0,
				Subtotal: 0,
			}
		}

		ticketSummary[classID].Quantity++
		ticketSummary[classID].Subtotal += price
		total += price
	}

	summaryList := make([]model.ClaimSessionTicketPricesResponse, 0, len(ticketSummary))
	for _, entry := range ticketSummary {
		summaryList = append(summaryList, *entry)
	}

	return summaryList, total
}
