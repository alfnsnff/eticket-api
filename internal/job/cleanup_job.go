package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"eticket-api/internal/repository"
	tx "eticket-api/pkg/utils/helper"

	"gorm.io/gorm"
)

type TicketCleanupJob struct {
	DB                *gorm.DB
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
	BatchSize         int
}

func NewTicketCleanupJob(db *gorm.DB, ticket_repository *repository.TicketRepository, session_repository *repository.SessionRepository, batchSize int) *TicketCleanupJob {
	return &TicketCleanupJob{
		DB:                db,
		TicketRepository:  ticket_repository,
		SessionRepository: session_repository,
		BatchSize:         batchSize,
	}
}

func (j *TicketCleanupJob) Run(ctx context.Context) error {
	log.Println("Running expired session cleanup job...")

	now := time.Now()

	expiredSessions, err := j.SessionRepository.FindExpired(j.DB, now, j.BatchSize)
	if err != nil {
		log.Printf("Error finding expired sessions: %v", err)
		return fmt.Errorf("failed to find expired sessions: %w", err)
	}

	if len(expiredSessions) == 0 {
		log.Println("No expired sessions found.")
		return nil // Nothing to do
	}

	log.Printf("Found %d expired sessions to process.", len(expiredSessions))

	for _, session := range expiredSessions {
		sessionErr := tx.Execute(ctx, j.DB, func(txDB *gorm.DB) error {

			ticketsToCancel, err := j.TicketRepository.FindManyBySessionID(txDB, session.ID) // Use txDB
			if err != nil {
				// Log this error, but return it to rollback this session's transaction
				return fmt.Errorf("failed to find tickets for expired session %d within transaction: %w", session.ID, err)
			}

			if len(ticketsToCancel) == 0 {
				log.Printf("Warning: Expired session %d has no linked tickets.", session.ID)
			} else {
				err = j.TicketRepository.CancelManyBySessionID(txDB, session.ID) // Assuming CancelManyBySessionID exists and accepts db *gorm.DB and sessionID
				if err != nil {
					return fmt.Errorf("failed to cancel tickets for session %d within transaction: %w", session.ID, err)
				}
				log.Printf("Cancelled %d tickets for expired session %d.", len(ticketsToCancel), session.ID)
			}

			err = j.SessionRepository.Delete(txDB, session) // Use txDB
			if err != nil {
				// Log this error, but return it to rollback this session's transaction
				return fmt.Errorf("failed to delete expired session %d within transaction: %w", session.ID, err)
			}
			log.Printf("Deleted expired session %d.", session.ID)

			return nil
		})

		if sessionErr != nil {
			log.Printf("Error cleaning up expired session %d: %v", session.ID, sessionErr)
		}
	}

	log.Println("Expired session cleanup job finished.")
	return nil
}
