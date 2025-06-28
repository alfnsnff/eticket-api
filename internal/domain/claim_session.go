package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type ClaimSession struct {
	ID         uint      `gorm:"column:id;primaryKey" json:"id"`
	SessionID  string    `gorm:"column:session_id;type:uuid;unique;not null"`
	ScheduleID uint      `gorm:"column:schedule_id;not null;index"`
	Status     string    `gorm:"column:status;type:varchar(24);not null"` //
	ExpiresAt  time.Time `gorm:"column:expires_at;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;not null"`

	Schedule   Schedule    `gorm:"foreignKey:ScheduleID" json:"schedule"` // Gorm will create the relationship
	ClaimItems []ClaimItem `gorm:"foreignKey:ClaimSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (cs *ClaimSession) TableName() string {
	return "claim_session"
}

type ClaimSessionRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *ClaimSession) error
	InsertBulk(ctx context.Context, conn gotann.Connection, sessions []*ClaimSession) error
	Update(ctx context.Context, conn gotann.Connection, entity *ClaimSession) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, sessions []*ClaimSession) error
	Delete(ctx context.Context, conn gotann.Connection, entity *ClaimSession) error
	DeleteBulk(ctx context.Context, conn gotann.Connection, entity []*ClaimSession) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*ClaimSession, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*ClaimSession, error)
	FindExpired(ctx context.Context, conn gotann.Connection, limit int) ([]*ClaimSession, error)
	FindByScheduleID(ctx context.Context, conn gotann.Connection, scheduleID uint) ([]*ClaimSession, error)
	FindBySessionID(ctx context.Context, conn gotann.Connection, uuid string) (*ClaimSession, error)
}
