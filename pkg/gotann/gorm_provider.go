package gotann

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// GormProvider holds the base *gorm.DB and config
type GormProvider struct {
	db     *gorm.DB
	config GormConfig
}

type GormConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	EnableLogging   bool
}

// GormConnection is a wrapper for *gorm.DB to implement Connection (non-transactional)
type GormConnection struct {
	db *gorm.DB
}

func NewGormConnection(db *gorm.DB) *GormConnection {
	return &GormConnection{db: db}
}

// GormTransaction wraps a GORM transaction and implements Connection + Transaction
type GormTransaction struct {
	tx         *gorm.DB
	id         string
	ctx        context.Context
	startTime  time.Time
	active     bool
	savepoints map[string]bool
	mu         sync.RWMutex
}

// --- GormProvider methods ---

func NewGormProvider(db *gorm.DB, config GormConfig) *GormProvider {
	return &GormProvider{
		db:     db,
		config: config,
	}
}

// Begin a transaction, returning GormTransaction
func (gp *GormProvider) Begin(ctx context.Context, opts TxOptions) (Transaction, error) {
	tx := gp.db.WithContext(ctx)

	// Isolation level and read-only options
	if opts.IsolationLevel != IsolationDefault {
		if level, err := gp.mapIsolationLevel(opts.IsolationLevel); err == nil {
			tx = tx.Set("gorm:isolation_level", level)
		}
	}
	if opts.ReadOnly {
		tx = tx.Set("gorm:read_only", true)
	}

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

// --- GormConnection implements Connection (forward all methods) ---

func (g *GormConnection) Create(value interface{}) *gorm.DB { return g.db.Create(value) }
func (g *GormConnection) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.First(dest, conds...)
}
func (g *GormConnection) FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.FirstOrCreate(dest, conds...)
}
func (g *GormConnection) FirstOrInit(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.FirstOrInit(dest, conds...)
}
func (g *GormConnection) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Find(dest, conds...)
}
func (g *GormConnection) Take(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Take(dest, conds...)
}
func (g *GormConnection) Last(dest interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Last(dest, conds...)
}
func (g *GormConnection) Save(value interface{}) *gorm.DB { return g.db.Save(value) }
func (g *GormConnection) Update(column string, value interface{}) *gorm.DB {
	return g.db.Update(column, value)
}
func (g *GormConnection) Updates(values interface{}) *gorm.DB { return g.db.Updates(values) }
func (g *GormConnection) UpdateColumn(column string, value interface{}) *gorm.DB {
	return g.db.UpdateColumn(column, value)
}
func (g *GormConnection) UpdateColumns(values interface{}) *gorm.DB {
	return g.db.UpdateColumns(values)
}
func (g *GormConnection) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return g.db.Delete(value, conds...)
}

func (g *GormConnection) Where(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Where(query, args...)
}
func (g *GormConnection) Not(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Not(query, args...)
}
func (g *GormConnection) Or(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Or(query, args...)
}
func (g *GormConnection) Select(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Select(query, args...)
}
func (g *GormConnection) Omit(columns ...string) *gorm.DB { return g.db.Omit(columns...) }
func (g *GormConnection) Joins(query string, args ...interface{}) *gorm.DB {
	return g.db.Joins(query, args...)
}
func (g *GormConnection) Preload(query string, args ...interface{}) *gorm.DB {
	return g.db.Preload(query, args...)
}
func (g *GormConnection) Group(name string) *gorm.DB { return g.db.Group(name) }
func (g *GormConnection) Having(query interface{}, args ...interface{}) *gorm.DB {
	return g.db.Having(query, args...)
}
func (g *GormConnection) Order(value interface{}) *gorm.DB      { return g.db.Order(value) }
func (g *GormConnection) Limit(limit int) *gorm.DB              { return g.db.Limit(limit) }
func (g *GormConnection) Offset(offset int) *gorm.DB            { return g.db.Offset(offset) }
func (g *GormConnection) Distinct(args ...interface{}) *gorm.DB { return g.db.Distinct(args...) }
func (g *GormConnection) Table(name string, args ...interface{}) *gorm.DB {
	return g.db.Table(name, args...)
}
func (g *GormConnection) Model(value interface{}) *gorm.DB { return g.db.Model(value) }
func (g *GormConnection) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return g.db.Scopes(funcs...)
}
func (g *GormConnection) Unscoped() *gorm.DB                   { return g.db.Unscoped() }
func (g *GormConnection) Attrs(attrs ...interface{}) *gorm.DB  { return g.db.Attrs(attrs...) }
func (g *GormConnection) Assign(attrs ...interface{}) *gorm.DB { return g.db.Assign(attrs...) }
func (g *GormConnection) Count(count *int64) *gorm.DB          { return g.db.Count(count) }

func (g *GormConnection) Raw(sql string, values ...interface{}) *gorm.DB {
	return g.db.Raw(sql, values...)
}
func (g *GormConnection) Exec(sql string, values ...interface{}) *gorm.DB {
	return g.db.Exec(sql, values...)
}
func (g *GormConnection) Scan(dest interface{}) *gorm.DB { return g.db.Scan(dest) }
func (g *GormConnection) Pluck(column string, dest interface{}) *gorm.DB {
	return g.db.Pluck(column, dest)
}
func (g *GormConnection) Row() *sql.Row            { return g.db.Row() }
func (g *GormConnection) Rows() (*sql.Rows, error) { return g.db.Rows() }
func (g *GormConnection) ScanRows(rows *sql.Rows, dest interface{}) error {
	return g.db.ScanRows(rows, dest)
}
func (g *GormConnection) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return g.db.Transaction(fc, opts...)
}
func (g *GormConnection) Begin(opts ...*sql.TxOptions) *gorm.DB { return g.db.Begin(opts...) }
func (g *GormConnection) Commit() *gorm.DB                      { return g.db.Commit() }
func (g *GormConnection) Rollback() *gorm.DB                    { return g.db.Rollback() }
func (g *GormConnection) SavePoint(name string) *gorm.DB        { return g.db.SavePoint(name) }
func (g *GormConnection) RollbackTo(name string) *gorm.DB       { return g.db.RollbackTo(name) }

func (g *GormConnection) WithContext(ctx context.Context) *gorm.DB    { return g.db.WithContext(ctx) }
func (g *GormConnection) Session(config *gorm.Session) *gorm.DB       { return g.db.Session(config) }
func (g *GormConnection) Debug() *gorm.DB                             { return g.db.Debug() }
func (g *GormConnection) Set(name string, value interface{}) *gorm.DB { return g.db.Set(name, value) }
func (g *GormConnection) Get(name string) (interface{}, bool)         { return g.db.Get(name) }
func (g *GormConnection) InstanceSet(name string, value interface{}) *gorm.DB {
	return g.db.InstanceSet(name, value)
}
func (g *GormConnection) InstanceGet(name string) (interface{}, bool) { return g.db.InstanceGet(name) }
func (g *GormConnection) Logger() logger.Interface                    { return g.db.Logger }
func (g *GormConnection) Statement() *gorm.Statement                  { return g.db.Statement }
func (g *GormConnection) RowsAffected() int64                         { return g.db.RowsAffected }
func (g *GormConnection) Error() error                                { return g.db.Error }
func (g *GormConnection) AutoMigrate(dst ...interface{}) error        { return g.db.AutoMigrate(dst...) }
func (g *GormConnection) Migrator() gorm.Migrator                     { return g.db.Migrator() }
func (g *GormConnection) Clauses(conds ...clause.Expression) *gorm.DB { return g.db.Clauses(conds...) }
func (g *GormConnection) Association(column string) *gorm.Association {
	return g.db.Association(column)
}
func (g *GormConnection) NamingStrategy() schema.Namer { return g.db.NamingStrategy }
func (g *GormConnection) AddError(err error) error     { return g.db.AddError(err) }
func (g *GormConnection) Use(plugin gorm.Plugin) error { return g.db.Use(plugin) }
func (g *GormConnection) Name() string                 { return g.db.Name() }
func (g *GormConnection) Dialector() gorm.Dialector    { return g.db.Dialector }
func (g *GormConnection) Context() context.Context     { return g.db.Statement.Context }
func (g *GormConnection) DB() (*sql.DB, error)         { return g.db.DB() }

// --- GormTransaction implements Connection + Transaction ---

// All Connection methods, just forward to tx (same as GormConnection)
func (gt *GormTransaction) Create(value interface{}) *gorm.DB { return gt.tx.Create(value) }
func (gt *GormTransaction) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.First(dest, conds...)
}
func (gt *GormTransaction) FirstOrCreate(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.FirstOrCreate(dest, conds...)
}
func (gt *GormTransaction) FirstOrInit(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.FirstOrInit(dest, conds...)
}
func (gt *GormTransaction) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.Find(dest, conds...)
}
func (gt *GormTransaction) Take(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.Take(dest, conds...)
}
func (gt *GormTransaction) Last(dest interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.Last(dest, conds...)
}
func (gt *GormTransaction) Save(value interface{}) *gorm.DB { return gt.tx.Save(value) }
func (gt *GormTransaction) Update(column string, value interface{}) *gorm.DB {
	return gt.tx.Update(column, value)
}
func (gt *GormTransaction) Updates(values interface{}) *gorm.DB { return gt.tx.Updates(values) }
func (gt *GormTransaction) UpdateColumn(column string, value interface{}) *gorm.DB {
	return gt.tx.UpdateColumn(column, value)
}
func (gt *GormTransaction) UpdateColumns(values interface{}) *gorm.DB {
	return gt.tx.UpdateColumns(values)
}
func (gt *GormTransaction) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return gt.tx.Delete(value, conds...)
}

func (gt *GormTransaction) Where(query interface{}, args ...interface{}) *gorm.DB {
	return gt.tx.Where(query, args...)
}
func (gt *GormTransaction) Not(query interface{}, args ...interface{}) *gorm.DB {
	return gt.tx.Not(query, args...)
}
func (gt *GormTransaction) Or(query interface{}, args ...interface{}) *gorm.DB {
	return gt.tx.Or(query, args...)
}
func (gt *GormTransaction) Select(query interface{}, args ...interface{}) *gorm.DB {
	return gt.tx.Select(query, args...)
}
func (gt *GormTransaction) Omit(columns ...string) *gorm.DB { return gt.tx.Omit(columns...) }
func (gt *GormTransaction) Joins(query string, args ...interface{}) *gorm.DB {
	return gt.tx.Joins(query, args...)
}
func (gt *GormTransaction) Preload(query string, args ...interface{}) *gorm.DB {
	return gt.tx.Preload(query, args...)
}
func (gt *GormTransaction) Group(name string) *gorm.DB { return gt.tx.Group(name) }
func (gt *GormTransaction) Having(query interface{}, args ...interface{}) *gorm.DB {
	return gt.tx.Having(query, args...)
}
func (gt *GormTransaction) Order(value interface{}) *gorm.DB      { return gt.tx.Order(value) }
func (gt *GormTransaction) Limit(limit int) *gorm.DB              { return gt.tx.Limit(limit) }
func (gt *GormTransaction) Offset(offset int) *gorm.DB            { return gt.tx.Offset(offset) }
func (gt *GormTransaction) Distinct(args ...interface{}) *gorm.DB { return gt.tx.Distinct(args...) }
func (gt *GormTransaction) Table(name string, args ...interface{}) *gorm.DB {
	return gt.tx.Table(name, args...)
}
func (gt *GormTransaction) Model(value interface{}) *gorm.DB { return gt.tx.Model(value) }
func (gt *GormTransaction) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *gorm.DB {
	return gt.tx.Scopes(funcs...)
}
func (gt *GormTransaction) Unscoped() *gorm.DB                   { return gt.tx.Unscoped() }
func (gt *GormTransaction) Attrs(attrs ...interface{}) *gorm.DB  { return gt.tx.Attrs(attrs...) }
func (gt *GormTransaction) Assign(attrs ...interface{}) *gorm.DB { return gt.tx.Assign(attrs...) }
func (gt *GormTransaction) Count(count *int64) *gorm.DB          { return gt.tx.Count(count) }

func (gt *GormTransaction) Raw(sql string, values ...interface{}) *gorm.DB {
	return gt.tx.Raw(sql, values...)
}
func (gt *GormTransaction) Exec(sql string, values ...interface{}) *gorm.DB {
	return gt.tx.Exec(sql, values...)
}
func (gt *GormTransaction) Scan(dest interface{}) *gorm.DB { return gt.tx.Scan(dest) }
func (gt *GormTransaction) Pluck(column string, dest interface{}) *gorm.DB {
	return gt.tx.Pluck(column, dest)
}
func (gt *GormTransaction) Row() *sql.Row            { return gt.tx.Row() }
func (gt *GormTransaction) Rows() (*sql.Rows, error) { return gt.tx.Rows() }
func (gt *GormTransaction) ScanRows(rows *sql.Rows, dest interface{}) error {
	return gt.tx.ScanRows(rows, dest)
}
func (gt *GormTransaction) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return gt.tx.Transaction(fc, opts...)
}
func (gt *GormTransaction) Begin(opts ...*sql.TxOptions) *gorm.DB { return gt.tx.Begin(opts...) }
func (gt *GormTransaction) Commit() *gorm.DB                      { return gt.tx.Commit() }
func (gt *GormTransaction) Rollback() *gorm.DB                    { return gt.tx.Rollback() }
func (gt *GormTransaction) SavePoint(name string) *gorm.DB        { return gt.tx.SavePoint(name) }
func (gt *GormTransaction) RollbackTo(name string) *gorm.DB       { return gt.tx.RollbackTo(name) }

func (gt *GormTransaction) WithContext(ctx context.Context) *gorm.DB { return gt.tx.WithContext(ctx) }
func (gt *GormTransaction) Session(config *gorm.Session) *gorm.DB    { return gt.tx.Session(config) }
func (gt *GormTransaction) Unwrap() *gorm.DB                         { return gt.tx }
func (gt *GormTransaction) Debug() *gorm.DB                          { return gt.tx.Debug() }
func (gt *GormTransaction) Set(name string, value interface{}) *gorm.DB {
	return gt.tx.Set(name, value)
}
func (gt *GormTransaction) Get(name string) (interface{}, bool) { return gt.tx.Get(name) }
func (gt *GormTransaction) InstanceSet(name string, value interface{}) *gorm.DB {
	return gt.tx.InstanceSet(name, value)
}
func (gt *GormTransaction) InstanceGet(name string) (interface{}, bool) {
	return gt.tx.InstanceGet(name)
}
func (gt *GormTransaction) Logger() logger.Interface             { return gt.tx.Logger }
func (gt *GormTransaction) Statement() *gorm.Statement           { return gt.tx.Statement }
func (gt *GormTransaction) RowsAffected() int64                  { return gt.tx.RowsAffected }
func (gt *GormTransaction) Error() error                         { return gt.tx.Error }
func (gt *GormTransaction) AutoMigrate(dst ...interface{}) error { return gt.tx.AutoMigrate(dst...) }
func (gt *GormTransaction) Migrator() gorm.Migrator              { return gt.tx.Migrator() }
func (gt *GormTransaction) Clauses(conds ...clause.Expression) *gorm.DB {
	return gt.tx.Clauses(conds...)
}
func (gt *GormTransaction) Association(column string) *gorm.Association {
	return gt.tx.Association(column)
}
func (gt *GormTransaction) NamingStrategy() schema.Namer { return gt.tx.NamingStrategy }
func (gt *GormTransaction) AddError(err error) error     { return gt.tx.AddError(err) }
func (gt *GormTransaction) Use(plugin gorm.Plugin) error { return gt.tx.Use(plugin) }
func (gt *GormTransaction) Name() string                 { return gt.tx.Name() }
func (gt *GormTransaction) Dialector() gorm.Dialector    { return gt.tx.Dialector }
func (gt *GormTransaction) Context() context.Context     { return gt.tx.Statement.Context }
func (gt *GormTransaction) DB() (*sql.DB, error)         { return gt.tx.DB() }

// --- Transaction interface specific methods ---

func (gt *GormTransaction) CommitTx() error {
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

func (gt *GormTransaction) RollbackTx() error {
	gt.mu.Lock()
	defer gt.mu.Unlock()
	if !gt.active {
		return ErrTransactionNotActive
	}
	err := gt.tx.Rollback().Error
	gt.active = false
	return err
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
