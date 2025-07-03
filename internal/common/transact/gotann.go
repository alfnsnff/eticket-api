package transact

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"

	"gorm.io/gorm"
)

// TransactionManager wraps the gotann manager for your domain
type Gotann struct {
	manager *gotann.Manager
}

// NewTransactionManager creates a new transaction manager for your app
func NewTransactionManager(db *gorm.DB) *Gotann {
	// Create GORM provider
	provider := gotann.NewGormProvider(db, gotann.GormConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		EnableLogging:   true,
	})

	// Build manager with your preferences
	manager := gotann.NewBuilder(provider).
		WithMaxConcurrentTransactions(10).
		WithDefaultTimeout(30 * time.Second).
		WithMetrics(true).
		WithTracing(false).
		Build()

	return &Gotann{
		manager: manager,
	}
}

// Execute provides the simple interface for your use cases
func (tm *Gotann) Execute(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.Execute(ctx, fn)
}

// ExecuteWithRetry provides retry capability
func (tm *Gotann) ExecuteWithRetry(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.ExecuteWithOptions(ctx, gotann.WithRetry(3, 100*time.Millisecond), fn)
}

// ExecuteReadOnly provides read-only transactions
func (tm *Gotann) ExecuteReadOnly(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.ExecuteWithOptions(ctx, gotann.ReadOnly(), fn)
}

// ExecuteBatch provides batch operations
func (tm *Gotann) ExecuteBatch(ctx context.Context, operations []gotann.BatchOperation) error {
	return tm.manager.ExecuteBatch(ctx, operations)
}

// Stats returns transaction statistics
func (tm *Gotann) Stats() gotann.TxStats {
	return tm.manager.Stats()
}

// Close gracefully shuts down the manager
func (tm *Gotann) Close() error {
	return tm.manager.Close()
}
