package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
)

// Map ClaimSession domain to ReadClaimSessionResponse model
func ClaimSessionToResponse(session *domain.ClaimSession) *model.ReadClaimSessionResponse {
	prices, total := BuildPriceBreakdown(session.Tickets)
	details := BuildTicketBreakdown(session.Tickets)
	return &model.ReadClaimSessionResponse{
		ID:        session.ID,
		SessionID: session.SessionID,
		Status:    session.Status,
		Schedule: model.ClaimSessionSchedule{
			ID: session.Schedule.ID,
			Ship: model.ClaimSessionScheduleShip{
				ID:       session.Schedule.Ship.ID,
				ShipName: session.Schedule.Ship.ShipName,
			},
			DepartureHarbor: model.ClaimSessionScheduleHarbor{
				ID:         session.Schedule.DepartureHarbor.ID,
				HarborName: session.Schedule.DepartureHarbor.HarborName,
			},
			ArrivalHarbor: model.ClaimSessionScheduleHarbor{
				ID:         session.Schedule.ArrivalHarbor.ID,
				HarborName: session.Schedule.ArrivalHarbor.HarborName,
			},
			DepartureDatetime: session.Schedule.DepartureDatetime,
			ArrivalDatetime:   session.Schedule.ArrivalDatetime,
		},
		ExpiresAt:   session.ExpiresAt,
		Prices:      prices,
		Tickets:     details,
		TotalAmount: total,
		CreatedAt:   session.CreatedAt,
		UpdatedAt:   session.UpdatedAt,
	}
}

// Build ticket breakdown
func BuildTicketBreakdown(tickets []domain.Ticket) []model.ClaimSessionTicketDetailResponse {
	result := make([]model.ClaimSessionTicketDetailResponse, len(tickets))
	for i, ticket := range tickets {
		result[i] = model.ClaimSessionTicketDetailResponse{
			TicketID: ticket.ID,
			Class: model.ClaimSessionTicketClassItem{
				ID:        ticket.Class.ID,
				ClassName: ticket.Class.ClassName,
				Type:      ticket.Class.Type,
			},
			Type:  ticket.Type,
			Price: ticket.Price,
		}
	}
	return result
}

// Build price breakdown
func BuildPriceBreakdown(tickets []domain.Ticket) ([]model.ClaimSessionTicketPricesResponse, float64) {
	ticketSummary := make(map[uint]*model.ClaimSessionTicketPricesResponse)
	var total float64

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
