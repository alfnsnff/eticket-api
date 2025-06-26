package gotann

import (
	"context"
	"time"
)

// Connection represents any database connection (regular or transactional)
type Connection interface{}

// Transaction represents a database transaction with advanced capabilities
type Transaction interface {
	Connection
	Commit() error
	Rollback() error
	Context() context.Context
	ID() string
	StartTime() time.Time
	IsActive() bool
	SetSavepoint(name string) error
	RollbackToSavepoint(name string) error
	ReleaseSavepoint(name string) error
}

// TransactionManager provides advanced transaction management
type TransactionManager interface {
	// Simple execution (most common use case)
	Execute(ctx context.Context, fn func(tx Transaction) error) error

	// Advanced execution with options
	ExecuteWithOptions(ctx context.Context, opts TxOptions, fn func(tx Transaction) error) error

	// Manual transaction control
	Begin(ctx context.Context) (Transaction, error)
	BeginWithOptions(ctx context.Context, opts TxOptions) (Transaction, error)

	// Batch operations
	ExecuteBatch(ctx context.Context, operations []BatchOperation) error

	// Health and monitoring
	Stats() TxStats
	Close() error
}

// TxOptions provides advanced transaction configuration
type TxOptions struct {
	IsolationLevel  IsolationLevel
	Timeout         time.Duration
	ReadOnly        bool
	RetryPolicy     *RetryPolicy
	Hooks           *TxHooks
	Context         map[string]interface{}
	DisableLogging  bool
	SavepointPrefix string
}

// BatchOperation represents a single operation in a batch
type BatchOperation struct {
	Name      string
	Operation func(tx Transaction) error
	OnError   func(err error) error
	Retryable bool
}

// IsolationLevel represents transaction isolation levels
type IsolationLevel int

const (
	IsolationDefault IsolationLevel = iota
	IsolationReadUncommitted
	IsolationReadCommitted
	IsolationRepeatableRead
	IsolationSerializable
)

// RetryPolicy defines retry behavior for transactions
type RetryPolicy struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffStrategy BackoffStrategy
	RetryableErrors []string
}

// BackoffStrategy defines how delays increase between retries
type BackoffStrategy int

const (
	BackoffLinear BackoffStrategy = iota
	BackoffExponential
	BackoffFixed
)

// TxHooks provides lifecycle hooks for transactions
type TxHooks struct {
	BeforeBegin    func(ctx context.Context) error
	AfterBegin     func(ctx context.Context, tx Transaction) error
	BeforeCommit   func(ctx context.Context, tx Transaction) error
	AfterCommit    func(ctx context.Context, tx Transaction) error
	BeforeRollback func(ctx context.Context, tx Transaction) error
	AfterRollback  func(ctx context.Context, tx Transaction, err error) error
	OnPanic        func(ctx context.Context, tx Transaction, recovered interface{}) error
}

// TxStats provides transaction statistics and monitoring
type TxStats struct {
	ActiveTransactions     int64
	TotalTransactions      int64
	CommittedTransactions  int64
	RolledBackTransactions int64
	FailedTransactions     int64
	AverageExecutionTime   time.Duration
	LongestTransaction     time.Duration
	DeadlockCount          int64
	RetryCount             int64
}
