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

type ClaimSessionRepository struct{}

func NewClaimSessionRepository() *ClaimSessionRepository {
	return &ClaimSessionRepository{}
}

func (csr *ClaimSessionRepository) Count(db *gorm.DB) (int64, error) {
	var total int64
	result := db.Model(&domain.ClaimSession{}).Count(&total)
	return total, result.Error
}

// func (r *ClaimSessionRepository) CountActiveReservedQuantity(db *gorm.DB, scheduleID, classID uint) (int64, error) {
// 	var total int64

// 	result := db.
// 		Table("claim_session").
// 		Select("COALESCE(SUM(claim_item.quantity), 0)").
// 		Joins("JOIN claim_item ON claim_item.claim_session_id = claim_session.id").
// 		Where("claim_session.schedule_id = ? AND claim_item.class_id = ?", scheduleID, classID).
// 		Where(`
// 			claim_session.status = ? OR
// 			(claim_session.status IN ? AND claim_session.expires_at > ?)
// 		`,
// 			enum.ClaimSessionSuccess.String(),
// 			enum.GetPendingClaimSessionStatuses(),
// 			time.Now(),
// 		).
// 		Scan(&total)

// 	return total, result.Error
// }

func (csr *ClaimSessionRepository) Insert(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Create(claim_session)
	return result.Error
}

func (csr *ClaimSessionRepository) InsertBulk(db *gorm.DB, sessions []*domain.ClaimSession) error {
	result := db.Create(sessions)
	return result.Error
}

func (csr *ClaimSessionRepository) Update(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Save(claim_session)
	return result.Error
}

func (csr *ClaimSessionRepository) UpdateBulk(db *gorm.DB, sessions []*domain.ClaimSession) error {
	result := db.Save(&sessions)
	return result.Error
}

func (csr *ClaimSessionRepository) Delete(db *gorm.DB, claim_session *domain.ClaimSession) error {
	result := db.Select(clause.Associations).Delete(claim_session)
	return result.Error
}

func (csr *ClaimSessionRepository) FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*domain.ClaimSession, error) {
	sessions := []*domain.ClaimSession{}
	query := db.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Tickets").
		Preload("Tickets.Class").
		Preload("ClaimItems").
		Preload("ClaimItems.Class")
	if search != "" {
		search = "%" + search + "%"
		query = query.Where("session_id ILIKE ?", search)
	}
	if sort == "" {
		sort = "id asc"
	} else {
		sort = strings.Replace(sort, ":", " ", 1)
	}
	err := query.Order(sort).Limit(limit).Offset(offset).Find(&sessions).Error
	return sessions, err
}

func (csr *ClaimSessionRepository) FindByID(db *gorm.DB, id uint) (*domain.ClaimSession, error) {
	session := new(domain.ClaimSession)
	result := db.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Tickets").
		Preload("Tickets.Class").
		Preload("ClaimItems").
		Preload("ClaimItems.Class").
		First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

// FindByUUID retrieves a ClaimSession domain by its SessionUUID.
func (csr *ClaimSessionRepository) FindBySessionID(db *gorm.DB, uuid string) (*domain.ClaimSession, error) {
	session := new(domain.ClaimSession)
	result := db.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Tickets").
		Preload("Tickets.Class").
		Preload("ClaimItems").
		Preload("ClaimItems.Class").
		Where("session_id = ?", uuid).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

func (r *ClaimSessionRepository) FindByScheduleID(tx *gorm.DB, scheduleID uint) ([]*domain.ClaimSession, error) {
	sessions := []*domain.ClaimSession{}
	result := tx.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("Tickets").
		Preload("Tickets.Class").
		Preload("ClaimItems").
		Preload("ClaimItems.Class").
		Where("schedule_id = ? AND status = ?", scheduleID, enum.ClaimSessionPendingData).
		Where("expires_at > ?", time.Now()).
		Find(&sessions)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return sessions, result.Error
}

func (csr *ClaimSessionRepository) FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*domain.ClaimSession, error) {
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
