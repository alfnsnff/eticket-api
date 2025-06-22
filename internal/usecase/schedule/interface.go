package schedule

import (
	"eticket-api/internal/contracts"
)

// Use aliases for commonly used interfaces to keep code clean
type (
	ScheduleRepository   = contracts.ScheduleRepository
	SessionRepository    = contracts.SessionRepository
	ClassRepository      = contracts.ClassRepository
	FareRepository       = contracts.FareRepository
	TicketRepository     = contracts.TicketRepository
	ShipRepository       = contracts.ShipRepository
	ManifestRepository   = contracts.ManifestRepository
	AllocationRepository = contracts.AllocationRepository
)
