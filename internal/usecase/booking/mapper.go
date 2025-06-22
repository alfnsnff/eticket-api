package booking

import (
	"eticket-api/internal/entity"
	"eticket-api/internal/model"
)

// Map Booking entity to ReadBookingResponse model
func ToReadBookingResponse(booking *entity.Booking) *model.ReadBookingResponse {
	if booking == nil {
		return nil
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
		OrderID: *booking.OrderID,
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
			DepartureDatetime: *booking.Schedule.DepartureDatetime,
			ArrivalDatetime:   *booking.Schedule.ArrivalDatetime,
		},
		CustomerName:    booking.CustomerName,
		CustomerAge:     booking.CustomerAge,
		CUstomerGender:  booking.CustomerGender,
		IDType:          booking.IDType,
		IDNumber:        booking.IDNumber,
		PhoneNumber:     booking.PhoneNumber,
		Email:           booking.Email,
		Status:          "completed", // Example status
		ReferenceNumber: booking.ReferenceNumber,
		BookedAt:        booking.CreatedAt,
		CreatedAt:       booking.CreatedAt,
		UpdatedAt:       booking.UpdatedAt,
		Tickets:         tickets,
	}
}

// Map a slice of Allocation entities to ReadAllocationResponse models
func ToReadBookingResponses(bookings []*entity.Booking) []*model.ReadBookingResponse {
	responses := make([]*model.ReadBookingResponse, len(bookings))
	for i, booking := range bookings {
		responses[i] = ToReadBookingResponse(booking)
	}
	return responses
}

// Map WriteBookingRequest model to Booking entity
func FromWriteBookingRequest(request *model.WriteBookingRequest) *entity.Booking {
	return &entity.Booking{
		OrderID:         request.OrderID,
		ScheduleID:      request.ScheduleID,
		CustomerName:    request.CustomerName,
		CustomerAge:     request.CustomerAge,
		CustomerGender:  request.CustomerGender,
		Email:           request.Email,
		PhoneNumber:     request.PhoneNumber,
		IDType:          request.IDType,
		IDNumber:        request.IDNumber,
		ReferenceNumber: request.ReferenceNumber,
	}
}

// Map UpdateBookingRequest model to Booking entity
func FromUpdateBookingRequest(request *model.UpdateBookingRequest, booking *entity.Booking) {
	booking.ScheduleID = request.ScheduleID
	booking.CustomerName = request.CustomerName
	booking.CustomerAge = request.CustomerAge
	booking.CustomerGender = request.CustomerGender
	booking.Email = request.Email
	booking.PhoneNumber = request.PhoneNumber
	booking.IDType = request.IDType
	booking.IDNumber = request.IDNumber
	booking.ReferenceNumber = request.ReferenceNumber
}
