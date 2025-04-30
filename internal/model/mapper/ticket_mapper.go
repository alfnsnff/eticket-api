package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToTicketModel(ticket *entities.Ticket) *model.ReadTicketResponse {
	response := new(model.ReadTicketResponse)
	copier.Copy(&response, &ticket) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToTicketsModel(tickets []*entities.Ticket) []*model.ReadTicketResponse {
	responses := []*model.ReadTicketResponse{}
	for _, ticket := range tickets {
		responses = append(responses, ToTicketModel(ticket))
	}
	return responses
}

func ToTicketEntity(request *model.WriteTicketRequest) *entities.Ticket {
	ticket := new(entities.Ticket)
	copier.Copy(&ticket, &request)
	return ticket
}
