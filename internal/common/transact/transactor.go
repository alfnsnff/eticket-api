package transact

import (
	"context"
	"eticket-api/pkg/gotann"
	"time"

	"gorm.io/gorm"
)

// TransactionManager wraps the gotann manager for your domain
type Transactor struct {
	manager *gotann.Manager
}

// NewTransactionManager creates a new transaction manager for your app
func NewTransactionManager(db *gorm.DB) *Transactor {
	// Create GORM provider
	provider := gotann.NewGormProvider(db, gotann.GormConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: time.Hour,
		EnableLogging:   true,
	})

	// Build manager with your preferences
	manager := gotann.NewBuilder(provider).
		WithMaxConcurrentTransactions(50).
		WithDefaultTimeout(30 * time.Second).
		WithMetrics(true).
		WithTracing(false).
		Build()

	return &Transactor{
		manager: manager,
	}
}

// Execute provides the simple interface for your use cases
func (tm *Transactor) Execute(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.Execute(ctx, fn)
}

// ExecuteWithRetry provides retry capability
func (tm *Transactor) ExecuteWithRetry(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.ExecuteWithOptions(ctx, gotann.WithRetry(3, 100*time.Millisecond), fn)
}

// ExecuteReadOnly provides read-only transactions
func (tm *Transactor) ExecuteReadOnly(ctx context.Context, fn func(tx gotann.Transaction) error) error {
	return tm.manager.ExecuteWithOptions(ctx, gotann.ReadOnly(), fn)
}

// ExecuteBatch provides batch operations
func (tm *Transactor) ExecuteBatch(ctx context.Context, operations []gotann.BatchOperation) error {
	return tm.manager.ExecuteBatch(ctx, operations)
}

// Stats returns transaction statistics
func (tm *Transactor) Stats() gotann.TxStats {
	return tm.manager.Stats()
}

// Close gracefully shuts down the manager
func (tm *Transactor) Close() error {
	return tm.manager.Close()
}
