package schedule

import (
	"eticket-api/internal/contracts"
)

// If you need use case-specific methods, extend the base interface
type ExtendedTicketRepository interface {
	contracts.TicketRepository
	// Add schedule-specific methods here if needed
	// e.g., GetTicketsByScheduleAndDate(db *gorm.DB, scheduleID uint, date time.Time) ([]*entity.Ticket, error)
}

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
