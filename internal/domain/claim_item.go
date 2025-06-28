package domain

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"
)

type ClaimItem struct {
	ID             uint      `gorm:"column:id;primaryKey"`
	ClaimSessionID uint      `gorm:"column:claim_session_id;not null;index"` // Foreign key to ClaimSession
	ClassID        uint      `gorm:"column:class_id;not null;index"`         // Foreign key to Class
	Quantity       int       `gorm:"column:quantity;not null"`               // Number of tickets requested for this class
	Subtotal       float64   `gorm:"column:subtotal;not null"`               // Total price for this class (Quantity * Price)
	CreatedAt      time.Time `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time `gorm:"column:updated_at;not null"`

	Class Class `gorm:"foreignKey:ClassID"` // Gorm will create the relationship
}

func (ci *ClaimItem) TableName() string {
	return "claim_item"
}

type ClaimItemRepository interface {
	Count(ctx context.Context, conn gotann.Connection) (int64, error)
	Insert(ctx context.Context, conn gotann.Connection, entity *ClaimItem) error
	InsertBulk(ctx context.Context, conn gotann.Connection, ClaimItemes []*ClaimItem) error
	Update(ctx context.Context, conn gotann.Connection, entity *ClaimItem) error
	UpdateBulk(ctx context.Context, conn gotann.Connection, ClaimItemes []*ClaimItem) error
	Delete(ctx context.Context, conn gotann.Connection, entity *ClaimItem) error
	FindAll(ctx context.Context, conn gotann.Connection, limit, offset int, sort, search string) ([]*ClaimItem, error)
	FindByID(ctx context.Context, conn gotann.Connection, id uint) (*ClaimItem, error)
}
