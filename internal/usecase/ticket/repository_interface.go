package ticket

import (
	"eticket-api/internal/domain"
)

// Use type alias to point to the shared contract
type (
	TicketRepository   = domain.TicketRepository
	ScheduleRepository = domain.ScheduleRepository
	ManifestRepository = domain.ManifestRepository
	FareRepository     = domain.FareRepository
)
