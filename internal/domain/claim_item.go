package domain

import (
	"time"

	"gorm.io/gorm"
)

type ClaimItem struct {
	ID             uint
	ClaimSessionID uint
	ClassID        uint
	Quantity       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (ci *ClaimItem) TableName() string {
	return "claim_item"
}

type ClaimItemRepository interface {
	Count(db *gorm.DB) (int64, error)
	CountActiveReservedQuantity(db *gorm.DB, scheduleID, classID uint) (int64, error)
	Insert(db *gorm.DB, entity *ClaimItem) error
	InsertBulk(db *gorm.DB, ClaimItemes []*ClaimItem) error
	Update(db *gorm.DB, entity *ClaimItem) error
	UpdateBulk(db *gorm.DB, ClaimItemes []*ClaimItem) error
	Delete(db *gorm.DB, entity *ClaimItem) error
	FindAll(db *gorm.DB, limit, offset int, sort, search string) ([]*ClaimItem, error)
	FindByID(db *gorm.DB, id uint) (*ClaimItem, error)
}
