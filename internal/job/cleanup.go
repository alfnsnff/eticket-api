package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"eticket-api/internal/repository"

	"gorm.io/gorm"
)

type CleanupJob struct {
	DB                *gorm.DB // Assuming you have a DB field for the transaction manager
	TicketRepository  *repository.TicketRepository
	SessionRepository *repository.SessionRepository
}

func NewCleanupJob(db *gorm.DB, ticket_repository *repository.TicketRepository, session_repository *repository.SessionRepository) *CleanupJob {
	return &CleanupJob{
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
		if err := j.SessionRepository.Delete(tx, session); err != nil {
			return fmt.Errorf("failed to delete expired session %d within transaction: %w", session.ID, err)
		}
		log.Printf("Deleted expired session %d.", session.ID)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Expired session cleanup job finished.")
	return nil
}
