package domain

import (
	"time"

	"gorm.io/gorm"
)

type ClaimItem struct {
	ID             uint      `gorm:"column:id;primaryKey"`
	ClaimSessionID uint      `gorm:"column:claim_session_id;not null;index"` // Foreign key to ClaimSession
	ClassID        uint      `gorm:"column:class_id;not null;index"`         // Foreign key to Class
	Quantity       int       `gorm:"column:quantity;not null"`               // Number of tickets requested for this class
	CreatedAt      time.Time `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"` // Gorm will create the relationship
}

func (ci *ClaimItem) TableName() string {
	return "claim_item"
}

type ClaimItemRepository interface {
	Count(db *gorm.DB) (int64, error)
	Insert(db *gorm.DB, entity *ClaimItem) error
	InsertBulk(db *gorm.DB, ClaimItemes []*ClaimItem) error
	Update(db *gorm.DB, entity *ClaimItem) error
	UpdateBulk(db *gorm.DB, ClaimItemes []*ClaimItem) error
	Delete(db *gorm.DB, entity *ClaimItem) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*ClaimItem, error)
	FindByID(db *gorm.DB, id uint) (*ClaimItem, error)
}
