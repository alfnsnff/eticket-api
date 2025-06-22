package repository

import (
	"errors"
	"strings"

	enum "eticket-api/internal/common/enums"
	"eticket-api/internal/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (ar *SessionRepository) Create(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Create(claim_session)
	return result.Error
}

func (ar *SessionRepository) Update(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Save(claim_session)
	return result.Error
}

func (ar *SessionRepository) Delete(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Select(clause.Associations).Delete(claim_session)
	return result.Error
}

func (csr *SessionRepository) Count(db *gorm.DB) (int64, error) {
	sessions := []*domain.ClaimSession{}
	var total int64
	result := db.Find(&sessions).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (csr *SessionRepository) GetAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.ClaimSession, error) {
	sessions := []*domain.ClaimSession{}

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

func (csr *SessionRepository) GetByID(db *gorm.DB, id uint) (*domain.ClaimSession, error) {
	session := new(domain.ClaimSession)
	result := db.Preload("Schedule").Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").Preload("Schedule.Ship").First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

// GetByUUID retrieves a ClaimSession domain by its SessionUUID.
func (csr *SessionRepository) GetByUUID(db *gorm.DB, uuid string) (*domain.ClaimSession, error) {
	var session domain.ClaimSession
	// Use the provided db instance (txDB from the use case)
	result := db.Preload("Schedule").Preload("Schedule.Route").
		Preload("Schedule.Route.DepartureHarbor").
		Preload("Schedule.Route.ArrivalHarbor").Preload("Schedule.Ship").Where("session_id = ?", uuid).First(&session)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil domain and nil error if not found
		}
		return nil, fmt.Errorf("failed to get claim session by UUID %s: %w", uuid, result.Error)
	}

	return &session, nil // Return pointer to the found domain
}

// GetByUUIDWithLock retrieves a ClaimSession domain by its SessionUUID with a lock.
func (csr *SessionRepository) GetByUUIDWithLock(db *gorm.DB, uuid string, forUpdate bool) (*domain.ClaimSession, error) {
	var session domain.ClaimSession
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
			return nil, nil // Return nil domain and nil error if not found
		}
		return nil, fmt.Errorf("failed to get claim session by ID %s with lock: %w", uuid, result.Error)
	}

	return &session, nil // Return pointer to the found domain
}

func (csr *SessionRepository) FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*domain.ClaimSession, error) {
	var sessions []*domain.ClaimSession

	result := db.Where(
		"(expires_at <= ? AND status NOT IN ?) OR status IN ?",
		expiryTime,
		enum.GetSuccessClaimSessionStatuses(),
		enum.GetFailedClaimSessionStatuses(),
	).Limit(limit).Find(&sessions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find expired claim sessions: %w", result.Error)
	}

	return sessions, nil
}
