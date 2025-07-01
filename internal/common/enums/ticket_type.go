package enum

// ShipStatus represents the status of a claim session
type TicketType int

const (
	TicketPassenger TicketType = iota
	TicketVehicle
)

func (ss TicketType) String() string {
	switch ss {
	case TicketPassenger:
		return "PASSENGER"
	case TicketVehicle:
		return "VEHICLE"

	default:
		return "UNKNOWN"
	}
}
