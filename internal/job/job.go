package job

import (
	"eticket-api/internal/injector"
	"eticket-api/internal/repository"
)

func Setup(ic *injector.Container) *CleanupJob {
	TicketRepository := repository.NewTicketRepository()
	SessionRepository := repository.NewSessionRepository()
	return NewCleanupJob(ic.Tx, TicketRepository, SessionRepository, 100)
}
