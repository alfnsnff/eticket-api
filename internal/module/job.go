package module

import (
	"eticket-api/internal/common/tx"
	"eticket-api/internal/job"
)

type JobModule struct {
	CleanupJob *job.CleanupJob
}

func NewJobModule(tx *tx.TxManager, repository *RepositoryModule) *JobModule {
	return &JobModule{
		CleanupJob: job.NewCleanupJob(tx, repository.TicketRepository, repository.SessionRepository),
	}
}
