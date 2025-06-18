package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"eticket-api/internal/common/tx"
	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type CleanupJob struct {
	Tx                *tx.TxManager
	DB                *gorm.DB // Assuming you have a DB field for the transaction manager
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewCleanupJob(tx *tx.TxManager, db *gorm.DB, ticket_repository *repository.TicketRepository, session_repository *repository.SessionRepository) *CleanupJob {
	return &CleanupJob{
		Tx:                tx,
		DB:                db,
		TicketRepository:  ticket_repository,
		SessionRepository: session_repository,
	}
}

func (j *CleanupJob) Run(ctx context.Context) error {
	tx := j.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			tx.Rollback()
		}
	}()

	log.Println("Running expired session cleanup job...")

	expiredSessions, err := j.SessionRepository.FindExpired(tx, time.Now(), 100)
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

		if err := j.SessionRepository.Delete(tx, session); err != nil {
			// Log this error, but return it to rollback this session's transaction
			return fmt.Errorf("failed to delete expired session %d within transaction: %w", session.ID, err)
		}
		log.Printf("Deleted expired session %d.", session.ID)

		return nil
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Expired session cleanup job finished.")
	return nil
}
