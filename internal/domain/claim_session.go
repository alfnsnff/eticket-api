package domain

import (
	"time"

	"gorm.io/gorm"
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
	Tickets    []Ticket    `gorm:"foreignKey:ClaimSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"tickets"`
	ClaimItems []ClaimItem `gorm:"foreignKey:ClaimSessionID" json:"claim_items"`
}

func (cs *ClaimSession) TableName() string {
	return "claim_session"
}

type ClaimSessionRepository interface {
	Count(db *gorm.DB) (int64, error)
	// CountActiveReservedQuantity(db *gorm.DB, scheduleID, classID uint) (int64, error)
	Insert(db *gorm.DB, entity *ClaimSession) error
	InsertBulk(db *gorm.DB, sessions []*ClaimSession) error
	Update(db *gorm.DB, entity *ClaimSession) error
	UpdateBulk(db *gorm.DB, sessions []*ClaimSession) error
	Delete(db *gorm.DB, entity *ClaimSession) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*ClaimSession, error)
	FindByID(db *gorm.DB, id uint) (*ClaimSession, error)
	FindExpired(db *gorm.DB, expiryTime time.Time, limit int) ([]*ClaimSession, error)
	FindByScheduleID(db *gorm.DB, scheduleID uint) ([]*ClaimSession, error)
	FindBySessionID(db *gorm.DB, uuid string) (*ClaimSession, error)
}
