package tx

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Execute runs a function inside a transaction with automatic rollback/commit
func Execute(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // rethrow panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
