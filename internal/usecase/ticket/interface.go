package ticket

import (
	"eticket-api/internal/contracts"
)

// Use type alias to point to the shared contract
type (
	TicketRepository   = contracts.TicketRepository
	ScheduleRepository = contracts.ScheduleRepository
	ManifestRepository = contracts.ManifestRepository
	FareRepository     = contracts.FareRepository
)
