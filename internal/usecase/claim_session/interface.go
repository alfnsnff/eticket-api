package claim_session

import (
	"eticket-api/internal/contracts"
)

// Use aliases for shared interfaces - eliminates duplication
type (
	SessionRepository    = contracts.SessionRepository
	BookingRepository    = contracts.BookingRepository
	ClassRepository      = contracts.ClassRepository
	FareRepository       = contracts.FareRepository
	TicketRepository     = contracts.TicketRepository
	ShipRepository       = contracts.ShipRepository
	ScheduleRepository   = contracts.ScheduleRepository
	ManifestRepository   = contracts.ManifestRepository
	AllocationRepository = contracts.AllocationRepository
)

// Only define interfaces here if claim_session has unique requirements
// For example, if you need claim session-specific methods:
// type ExtendedSessionRepository interface {
//     contracts.SessionRepository
//     GetExpiredSessionsByUser(db *gorm.DB, userID uint) ([]*entity.ClaimSession, error)
// }
