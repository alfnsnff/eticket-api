package tx

import (
	"context"

	"gorm.io/gorm"
)

// TxManager abstracts transaction execution
type TxManager interface {
	Execute(ctx context.Context, fn func(tx *gorm.DB) error) error
}

// GormTxManager is a concrete implementation of TxManager using GORM
type GormTxManager struct {
	DB *gorm.DB
}

// NewGormTxManager creates a new GormTxManager
func NewGormTxManager(db *gorm.DB) *GormTxManager {
	return &GormTxManager{DB: db}
}

// Execute satisfies the TxManager interface using the global Execute logic
func (g *GormTxManager) Execute(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return Execute(ctx, g.DB, fn)
}
