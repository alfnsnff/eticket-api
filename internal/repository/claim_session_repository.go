package repository

import (
	"context"
	"errors"
	"strings"

	enum "eticket-api/internal/common/enums"
	"eticket-api/internal/domain"
	"eticket-api/pkg/gotann"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClaimSessionRepository struct {
	DB *gorm.DB
}

func NewClaimSessionRepository(db *gorm.DB) *ClaimSessionRepository {
	return &ClaimSessionRepository{DB: db}
}

func (r *ClaimSessionRepository) Count(ctx context.Context, conn gotann.Connection) (int64, error) {
	var total int64
	result := conn.Model(&domain.ClaimSession{}).Count(&total)
	return total, result.Error
}

func (r *ClaimSessionRepository) Insert(ctx context.Context, conn gotann.Connection, claim_session *domain.ClaimSession) error {
	result := conn.Create(claim_session)
	return result.Error
}

func (r *ClaimSessionRepository) InsertBulk(ctx context.Context, conn gotann.Connection, sessions []*domain.ClaimSession) error {
	result := conn.Create(sessions)
	return result.Error
}

func (r *ClaimSessionRepository) Update(ctx context.Context, conn gotann.Connection, claim_session *domain.ClaimSession) error {
	result := conn.Save(claim_session)
	return result.Error
}

func (r *ClaimSessionRepository) UpdateBulk(ctx context.Context, conn gotann.Connection, sessions []*domain.ClaimSession) error {
	result := conn.Save(&sessions)
	return result.Error
}

func (r *ClaimSessionRepository) Delete(ctx context.Context, conn gotann.Connection, claim_session *domain.ClaimSession) error {
	result := conn.Select(clause.Associations).Delete(claim_session)
	return result.Error
}

func (r *ClaimSessionRepository) DeleteBulk(ctx context.Context, conn gotann.Connection, claim_sessions []*domain.ClaimSession) error {
	result := conn.Select(clause.Associations).Delete(claim_sessions)
	return result.Error
}

func (r *ClaimSessionRepository) FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*domain.ClaimSession, error) {
	sessions := []*domain.ClaimSession{}
	query := conn.Model(&domain.ClaimSession{}).Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
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

func (r *ClaimSessionRepository) FindByID(ctx context.Context, conn gotann.Connection, id uint) (*domain.ClaimSession, error) {
	session := new(domain.ClaimSession)
	result := conn.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("ClaimItems").
		Preload("ClaimItems.Class").
		First(&session, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

// FindByUUID retrieves a ClaimSession domain by its SessionUUID.
func (r *ClaimSessionRepository) FindBySessionID(ctx context.Context, conn gotann.Connection, uuid string) (*domain.ClaimSession, error) {
	session := new(domain.ClaimSession)
	result := conn.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
		Preload("ClaimItems").
		Preload("ClaimItems.Class").
		Where("session_id = ?", uuid).First(&session)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return session, result.Error
}

func (r *ClaimSessionRepository) FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*domain.ClaimSession, error) {
	sessions := []*domain.ClaimSession{}
	result := conn.Preload("Schedule").
		Preload("Schedule.DepartureHarbor").
		Preload("Schedule.ArrivalHarbor").
		Preload("Schedule.Ship").
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

func (r *ClaimSessionRepository) FindExpired(ctx context.Context, conn gotann.Connection, limit int) ([]*domain.ClaimSession, error) {
	var sessions []*domain.ClaimSession
	now := time.Now()
	result := conn.Where(
		"(expires_at <= ? AND status NOT IN ?) OR status IN ?",
		now,
		enum.GetSuccessClaimSessionStatuses(),
		enum.GetFailedClaimSessionStatuses(),
	).Limit(limit).Find(&sessions)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find expired claim sessions: %w", result.Error)
	}

	return sessions, nil
}
