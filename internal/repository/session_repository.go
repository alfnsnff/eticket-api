package repository

import (
	"errors"
	"eticket-api/internal/domain/entity"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository struct {
	Repository[entity.ClaimSession]
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (shr *SessionRepository) GetAll(db *gorm.DB) ([]*entity.ClaimSession, error) {
	sessions := []*entity.ClaimSession{}
	result := db.Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return sessions, nil
}

func (shr *SessionRepository) GetByID(db *gorm.DB, id uint) (*entity.ClaimSession, error) {
	session := new(entity.ClaimSession)
	result := db.First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

// GetByUUID retrieves a ClaimSession entity by its SessionUUID.
func (r *SessionRepository) GetByUUID(db *gorm.DB, uuid string) (*entity.ClaimSession, error) {
	var session entity.ClaimSession
	// Use the provided db instance (txDB from the use case)
	result := db.Where("session_id = ?", uuid).First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil entity and nil error if not found
		}
		return nil, fmt.Errorf("failed to get claim session by UUID %s: %w", uuid, result.Error)
	}

	return &session, nil // Return pointer to the found entity
}

// GetByUUIDWithLock retrieves a ClaimSession entity by its SessionUUID with a lock.
func (r *SessionRepository) GetByUUIDWithLock(db *gorm.DB, uuid string, forUpdate bool) (*entity.ClaimSession, error) {
	var session entity.ClaimSession
	query := db.Where("session_id = ?", uuid)

	if forUpdate {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"}) // Add FOR UPDATE lock
	} else {
		// Add other lock types if needed, e.g., FOR SHARE
		// query = query.Clauses(clause.Locking{Strength: "SHARE"})
	}

	result := query.First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil entity and nil error if not found
		}
		return nil, fmt.Errorf("failed to get claim session by ID %s with lock: %w", uuid, result.Error)
	}

	return &session, nil // Return pointer to the found entity
}

// FindExpired retrieves ClaimSession entities whose ExpiresAt is in the past.
func (r *SessionRepository) FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*entity.ClaimSession, error) {
	var sessions []*entity.ClaimSession

	// Use the provided db instance (txDB from the job or a new session if not in tx)
	// Find sessions where expires_at is less than or equal to the provided expiryTime (usually time.Now())
	// Add a limit for batch processing in the cleanup job
	result := db.Where("expires_at <= ?", expiryTime).Limit(limit).Find(&sessions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find expired claim sessions: %w", result.Error)
	}

	// Return the slice of sessions (it will be empty if none were found) and a nil error
	return sessions, nil
}
