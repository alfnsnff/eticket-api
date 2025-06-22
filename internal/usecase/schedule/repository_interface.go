package schedule

import (
	"eticket-api/internal/domain"
)

// Use aliases for commonly used interfaces to keep code clean
type (
	ScheduleRepository   = domain.ScheduleRepository
	ClassRepository      = domain.ClassRepository
	FareRepository       = domain.FareRepository
	TicketRepository     = domain.TicketRepository
	ShipRepository       = domain.ShipRepository
	ManifestRepository   = domain.ManifestRepository
	AllocationRepository = domain.AllocationRepository
)
