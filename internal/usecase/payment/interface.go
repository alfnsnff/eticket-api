package payment

import (
	"eticket-api/internal/contracts"
)

// Use type aliases to point to the shared contracts
type (
	BookingRepository      = contracts.BookingRepository
	TicketRepository       = contracts.TicketRepository
	ClaimSessionRepository = contracts.ClaimSessionRepository
)
