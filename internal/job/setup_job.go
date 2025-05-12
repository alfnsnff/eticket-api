package job

import (
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

func SetupJob(db *gorm.DB) *TicketCleanupJob {
	TicketRepository := repository.NewTicketRepository()
	SessionRepository := repository.NewSessionRepository()
	return NewTicketCleanupJob(db, TicketRepository, SessionRepository, 100)
}
