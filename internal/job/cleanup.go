package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"eticket-api/internal/common/tx"
	"eticket-api/internal/entity"
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type CleanupJob struct {
	Tx                *tx.TxManager
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewCleanupJob(tx *tx.TxManager, ticket_repository *repository.TicketRepository, session_repository *repository.SessionRepository) *CleanupJob {
	return &CleanupJob{
		Tx:                tx,
		TicketRepository:  ticket_repository,
		SessionRepository: session_repository,
	}
}

func (j *CleanupJob) Run(ctx context.Context) error {
	log.Println("Running expired session cleanup job...")

	now := time.Now()
	expiredSessions := []*entity.ClaimSession{}

	j.Tx.Execute(ctx, func(tx *gorm.DB) error {
		var err error
		expiredSessions, err = j.SessionRepository.FindExpired(tx, now, 100)
		if err != nil {
			log.Printf("Error finding expired sessions: %v", err)
			return fmt.Errorf("failed to find expired sessions: %w", err)
		}

		return nil
	})

	if len(expiredSessions) == 0 {
		log.Println("No expired sessions found.")
		return nil // Nothing to do
	}

	log.Printf("Found %d expired sessions to process.", len(expiredSessions))

	for _, session := range expiredSessions {
		sessionErr := j.Tx.Execute(ctx, func(tx *gorm.DB) error {

			ticketsToCancel, err := j.TicketRepository.FindManyBySessionID(tx, session.ID) // Use tx
			if err != nil {
				// Log this error, but return it to rollback this session's transaction
				return fmt.Errorf("failed to find tickets for expired session %d within transaction: %w", session.ID, err)
			}

			if len(ticketsToCancel) == 0 {
				log.Printf("Warning: Expired session %d has no linked tickets.", session.ID)
			} else {
				err = j.TicketRepository.CancelManyBySessionID(tx, session.ID) // Assuming CancelManyBySessionID exists and accepts db *gorm.DB and sessionID
				if err != nil {
					return fmt.Errorf("failed to cancel tickets for session %d within transaction: %w", session.ID, err)
				}
				log.Printf("Cancelled %d tickets for expired session %d.", len(ticketsToCancel), session.ID)
			}

			err = j.SessionRepository.Delete(tx, session) // Use tx
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
