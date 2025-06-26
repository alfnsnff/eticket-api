package gotann

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormProvider implements the Provider interface for GORM
type GormProvider struct {
	db     *gorm.DB
	config GormConfig
	mu     sync.RWMutex
}

// GormConfig holds GORM-specific configuration
type GormConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	EnableLogging   bool
}

// GormTransaction wraps a GORM transaction
type GormTransaction struct {
	tx         *gorm.DB
	id         string
	ctx        context.Context
	startTime  time.Time
	active     bool
	savepoints map[string]bool
	mu         sync.RWMutex
}

// NewGormProvider creates a new GORM provider
func NewGormProvider(db *gorm.DB, config GormConfig) *GormProvider {
	return &GormProvider{
		db:     db,
		config: config,
	}
}

// Begin starts a new GORM transaction
func (gp *GormProvider) Begin(ctx context.Context, opts TxOptions) (Transaction, error) {
	tx := gp.db.WithContext(ctx)

	// Apply isolation level if supported
	if opts.IsolationLevel != IsolationDefault {
		if level, err := gp.mapIsolationLevel(opts.IsolationLevel); err == nil {
			tx = tx.Set("gorm:isolation_level", level)
		}
	}

	// Apply read-only mode
	if opts.ReadOnly {
		tx = tx.Set("gorm:read_only", true)
	}

	// Begin the transaction
	gormTx := tx.Begin()
	if gormTx.Error != nil {
		return nil, fmt.Errorf("failed to begin GORM transaction: %w", gormTx.Error)
	}

	return &GormTransaction{
		tx:         gormTx,
		id:         uuid.New().String(),
		ctx:        ctx,
		startTime:  time.Now(),
		active:     true,
		savepoints: make(map[string]bool),
	}, nil
}

// SupportsIsolationLevel checks if the isolation level is supported
func (gp *GormProvider) SupportsIsolationLevel(level IsolationLevel) bool {
	switch level {
	case IsolationDefault, IsolationReadCommitted, IsolationRepeatableRead, IsolationSerializable:
		return true
	case IsolationReadUncommitted:
		// Most databases support this, but not all
		return true
	default:
		return false
	}
}

// SupportsSavepoints returns true as GORM supports savepoints
func (gp *GormProvider) SupportsSavepoints() bool {
	return true
}

// MaxConnections returns the maximum number of connections
func (gp *GormProvider) MaxConnections() int {
	return gp.config.MaxOpenConns
}

// HealthCheck verifies the database connection
func (gp *GormProvider) HealthCheck() error {
	sqlDB, err := gp.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// mapIsolationLevel converts our isolation level to GORM format
func (gp *GormProvider) mapIsolationLevel(level IsolationLevel) (string, error) {
	switch level {
	case IsolationReadUncommitted:
		return "READ UNCOMMITTED", nil
	case IsolationReadCommitted:
		return "READ COMMITTED", nil
	case IsolationRepeatableRead:
		return "REPEATABLE READ", nil
	case IsolationSerializable:
		return "SERIALIZABLE", nil
	default:
		return "", fmt.Errorf("unsupported isolation level: %v", level)
	}
}

// GormTransaction implementation
func (gt *GormTransaction) Commit() error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	if !gt.active {
		return ErrTransactionNotActive
	}

	err := gt.tx.Commit().Error
	if err == nil {
		gt.active = false
	}
	return err
}

func (gt *GormTransaction) Rollback() error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	if !gt.active {
		return ErrTransactionNotActive
	}

	err := gt.tx.Rollback().Error
	gt.active = false
	return err
}

func (gt *GormTransaction) Context() context.Context {
	return gt.ctx
}

func (gt *GormTransaction) ID() string {
	return gt.id
}

func (gt *GormTransaction) StartTime() time.Time {
	return gt.startTime
}

func (gt *GormTransaction) IsActive() bool {
	gt.mu.RLock()
	defer gt.mu.RUnlock()
	return gt.active
}

func (gt *GormTransaction) SetSavepoint(name string) error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	if !gt.active {
		return ErrTransactionNotActive
	}

	err := gt.tx.Exec(fmt.Sprintf("SAVEPOINT %s", name)).Error
	if err == nil {
		gt.savepoints[name] = true
	}
	return err
}

func (gt *GormTransaction) RollbackToSavepoint(name string) error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	if !gt.active {
		return ErrTransactionNotActive
	}

	if !gt.savepoints[name] {
		return ErrSavepointNotFound
	}

	return gt.tx.Exec(fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name)).Error
}

func (gt *GormTransaction) ReleaseSavepoint(name string) error {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	if !gt.active {
		return ErrTransactionNotActive
	}

	if !gt.savepoints[name] {
		return ErrSavepointNotFound
	}

	err := gt.tx.Exec(fmt.Sprintf("RELEASE SAVEPOINT %s", name)).Error
	if err == nil {
		delete(gt.savepoints, name)
	}
	return err
}

// GetGormDB returns the underlying GORM transaction for repository use
func (gt *GormTransaction) GetGormDB() *gorm.DB {
	return gt.tx
}

// Convenience method to get GORM DB from any connection
func GetGormDB(conn Connection) (*gorm.DB, error) {
	switch c := conn.(type) {
	case *GormTransaction:
		return c.GetGormDB(), nil
	case *gorm.DB:
		return c, nil
	default:
		return nil, fmt.Errorf("connection is not a GORM connection")
	}
}
