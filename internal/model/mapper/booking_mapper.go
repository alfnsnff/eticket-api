package mapper

import (
	"eticket-api/internal/domain/entities"
	"eticket-api/internal/model"

	"github.com/jinzhu/copier"
)

func ToBookingModel(booking *entities.Booking) *model.ReadBookingResponse {
	// var response *model.ReadBookingResponse
	response := new(model.ReadBookingResponse)
	copier.Copy(&response, &booking) // Automatically maps matching fields
	return response
}

// Convert a slice of Ticket entities to DTO slice
func ToBookingsModel(bookings []*entities.Booking) []*model.ReadBookingResponse {
	responses := []*model.ReadBookingResponse{}
	for _, booking := range bookings {
		responses = append(responses, ToBookingModel(booking))
	}
	return responses
}

func ToBookingEntity(request *model.WriteBookingRequest) *entities.Booking {
	booking := new(entities.Booking)
	copier.Copy(&booking, &request)
	return booking
}

// // Convert BookingReq DTO to Booking entity
// func ToBookingEntity(request *model.CreateBookingTicketRequest) (entities.Booking, []entities.Ticket) {
// 	var booking entities.Booking
// 	var tickets []entities.Ticket

// 	// Automatically copy matching fields
// 	copier.Copy(&booking, &bookingCreate)

// 	// Convert tickets manually (copier doesn't automatically handle slices of nested structs)
// 	for _, t := range bookingCreate.Tickets {
// 		var ticket entities.Ticket
// 		copier.Copy(&ticket, &t)
// 		tickets = append(tickets, ticket)
// 	}

// 	return booking, tickets
// }
