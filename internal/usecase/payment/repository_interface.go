package payment

import (
	"eticket-api/internal/domain"
)

// Use type aliases to point to the shared domain
type (
	BookingRepository      = domain.BookingRepository
	TicketRepository       = domain.TicketRepository
	ClaimSessionRepository = domain.ClaimSessionRepository
)
