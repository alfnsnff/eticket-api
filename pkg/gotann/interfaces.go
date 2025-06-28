package gotann

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Connection represents any database connection (regular or transactional) that exposes a GORM-like API
type Connection interface {
	// Basic CRUD
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB
	FirstOrInit(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Take(dest interface{}, conds ...interface{}) *gorm.DB
	Last(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Updates(values interface{}) *gorm.DB
	UpdateColumn(column string, value interface{}) *gorm.DB
	UpdateColumns(values interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB

	// Query
	Where(query interface{}, args ...interface{}) *gorm.DB
	Not(query interface{}, args ...interface{}) *gorm.DB
	Or(query interface{}, args ...interface{}) *gorm.DB
	Select(query interface{}, args ...interface{}) *gorm.DB
	Omit(columns ...string) *gorm.DB
	Joins(query string, args ...interface{}) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
	Group(name string) *gorm.DB
	Having(query interface{}, args ...interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
	Limit(limit int) *gorm.DB
	Offset(offset int) *gorm.DB
	Distinct(args ...interface{}) *gorm.DB
	Table(name string, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB
	Unscoped() *gorm.DB
	Attrs(attrs ...interface{}) *gorm.DB
	Assign(attrs ...interface{}) *gorm.DB
	Count(count *int64) *gorm.DB

	// Advanced
	Raw(sql string, values ...interface{}) *gorm.DB
	Exec(sql string, values ...interface{}) *gorm.DB
	Scan(dest interface{}) *gorm.DB
	Pluck(column string, dest interface{}) *gorm.DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Commit() *gorm.DB
	Rollback() *gorm.DB
	SavePoint(name string) *gorm.DB
	RollbackTo(name string) *gorm.DB

	// Context and Callbacks
	WithContext(ctx context.Context) *gorm.DB
	Session(config *gorm.Session) *gorm.DB
	Unwrap() *gorm.DB
	Debug() *gorm.DB
	Set(name string, value interface{}) *gorm.DB
	Get(name string) (interface{}, bool)
	InstanceSet(name string, value interface{}) *gorm.DB
	InstanceGet(name string) (interface{}, bool)

	// Logger
	Logger() logger.Interface

	// Statement and Schema
	Statement() *gorm.Statement
	RowsAffected() int64
	Error() error

	// Migrator
	AutoMigrate(dst ...interface{}) error
	Migrator() gorm.Migrator

	// Clause
	Clauses(conds ...clause.Expression) *gorm.DB

	// Association
	Association(column string) *gorm.Association

	// NamingStrategy
	NamingStrategy() schema.Namer

	// Utilities
	AddError(err error) error
	Use(plugin gorm.Plugin) error

	// Utility
	Name() string
	Dialector() gorm.Dialector

	// Context
	Context() context.Context

	// DB
	DB() (*sql.DB, error)
}

// Transaction represents a database transaction with advanced capabilities
type Transaction interface {
	Connection
	CommitTx() error
	RollbackTx() error
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
