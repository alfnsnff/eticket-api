package mapper

import (
	"eticket-api/internal/domain"
	"eticket-api/internal/model"
	"fmt"
)

func TicketToItem(ticket *domain.Ticket) model.OrderItem {
	name := "Tiket " + ticket.Type
	if ticket.Type == "passenger" && ticket.PassengerName != nil {
		name = fmt.Sprintf("Tiket Penumpang - %s", *ticket.PassengerName)
	}
	if ticket.Type == "vehicle" && ticket.LicensePlate != nil {
		name = fmt.Sprintf("Tiket Kendaraan - %s", *ticket.LicensePlate)
	}
	return model.OrderItem{
		SKU:      ticket.Class.ClassName,
		Name:     name,              // e.g. "Passenger", "Vehicle"
		Price:    int(ticket.Price), // pastikan tipe t.Price sesuai (float32 ke int)
		Quantity: 1,                 // satu tiket per entri
	}
}
