package repository

import (
	"errors"
	"strings"

	"eticket-api/internal/entity"
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

func (csr *SessionRepository) Count(db *gorm.DB) (int64, error) {
	sessions := []*entity.ClaimSession{}
	var total int64
	result := db.Find(&sessions).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (csr *SessionRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*entity.ClaimSession, error) {
	sessions := []*entity.ClaimSession{}

	query := db.Preload("Schedule").
		Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").
		Preload("Schedule.Ship")

	if search != "" {
		search = "%" + search + "%"
		query = query.Where("session_id ILIKE ?", search)
	}

	// 🔃 Sort (with default)
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}

	err := query.Order(sort).Limit(limit).Offset(offset).Find(&sessions).Error
	return sessions, err
}

func (csr *SessionRepository) GetByID(db *gorm.DB, id uint) (*entity.ClaimSession, error) {
	session := new(entity.ClaimSession)
	result := db.Preload("Schedule").Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").Preload("Schedule.Ship").First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

// GetByUUID retrieves a ClaimSession entity by its SessionUUID.
func (csr *SessionRepository) GetByUUID(db *gorm.DB, uuid string) (*entity.ClaimSession, error) {
	var session entity.ClaimSession
	// Use the provided db instance (txDB from the use case)
	result := db.Preload("Schedule").Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").Preload("Schedule.Ship").Where("session_id = ?", uuid).First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil entity and nil error if not found
		}
		return nil, fmt.Errorf("failed to get claim session by UUID %s: %w", uuid, result.Error)
	}

	return &session, nil // Return pointer to the found entity
}

// GetByUUIDWithLock retrieves a ClaimSession entity by its SessionUUID with a lock.
func (csr *SessionRepository) GetByUUIDWithLock(db *gorm.DB, uuid string, forUpdate bool) (*entity.ClaimSession, error) {
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

func (csr *SessionRepository) FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*entity.ClaimSession, error) {
	var sessions []*entity.ClaimSession

	result := db.Where("expires_at <= ?", expiryTime).Limit(limit).Find(&sessions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find expired claim sessions: %w", result.Error)
	}

	return sessions, nil
}
