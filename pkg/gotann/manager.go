package gotann

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Manager is the core transaction manager implementation
type Manager struct {
	provider      Provider
	config        Config
	stats         *stats
	logger        *slog.Logger
	mu            sync.RWMutex
	activeTxs     map[string]*transactionContext
	closed        bool
	shutdownCh    chan struct{}
	cleanupTicker *time.Ticker
}

// Provider interface abstracts the underlying database
type Provider interface {
	Begin(ctx context.Context, opts TxOptions) (Transaction, error)
	SupportsIsolationLevel(level IsolationLevel) bool
	SupportsSavepoints() bool
	MaxConnections() int
	HealthCheck() error
}

// Config holds manager configuration
type Config struct {
	MaxConcurrentTx    int
	DefaultTimeout     time.Duration
	CleanupInterval    time.Duration
	EnableMetrics      bool
	EnableTracing      bool
	LogLevel           slog.Level
	DeadlockRetryDelay time.Duration
	MaxRetryAttempts   int
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		MaxConcurrentTx:    100,
		DefaultTimeout:     30 * time.Second,
		CleanupInterval:    1 * time.Minute,
		EnableMetrics:      true,
		EnableTracing:      false,
		LogLevel:           slog.LevelInfo,
		DeadlockRetryDelay: 100 * time.Millisecond,
		MaxRetryAttempts:   3,
	}
}

// NewManager creates a new transaction manager
func NewManager(provider Provider, config Config, logger *slog.Logger) *Manager {
	if logger == nil {
		logger = slog.Default()
	}

	m := &Manager{
		provider:      provider,
		config:        config,
		stats:         newStats(),
		logger:        logger,
		activeTxs:     make(map[string]*transactionContext),
		shutdownCh:    make(chan struct{}),
		cleanupTicker: time.NewTicker(config.CleanupInterval),
	}

	// Start background cleanup goroutine
	go m.cleanupLoop()

	return m
}

// Execute is the main entry point - super simple to use!
func (m *Manager) Execute(ctx context.Context, fn func(tx Transaction) error) error {
	return m.ExecuteWithOptions(ctx, TxOptions{}, fn)
}

// ExecuteWithOptions provides advanced execution with full control
func (m *Manager) ExecuteWithOptions(ctx context.Context, opts TxOptions, fn func(tx Transaction) error) error {
	if m.closed {
		return ErrManagerClosed
	}

	// Apply defaults
	opts = m.applyDefaults(opts)

	// Apply timeout context
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Execute with retry policy
	if opts.RetryPolicy != nil {
		return m.executeWithRetry(ctx, opts, fn)
	}

	return m.executeOnce(ctx, opts, fn)
}
func (m *Manager) executeOnce(ctx context.Context, opts TxOptions, fn func(tx Transaction) error) (err error) {
	// Check concurrent transaction limit
	if !m.canStartTransaction() {
		return ErrTooManyTransactions
	}

	// Execute hooks
	if opts.Hooks != nil && opts.Hooks.BeforeBegin != nil {
		if err := opts.Hooks.BeforeBegin(ctx); err != nil {
			return fmt.Errorf("before begin hook failed: %w", err)
		}
	}

	// Begin transaction
	tx, err := m.Begin(ctx)
	if err != nil {
		atomic.AddInt64(&m.stats.failedTx, 1)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Track transaction
	txCtx := &transactionContext{
		tx:        tx,
		startTime: time.Now(),
		context:   ctx,
		options:   opts,
	}

	m.addActiveTransaction(tx.ID(), txCtx)
	defer m.removeActiveTransaction(tx.ID())

	// Setup panic recovery
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Transaction panic occurred",
				"tx_id", tx.ID(),
				"panic", r,
			)

			if opts.Hooks != nil && opts.Hooks.OnPanic != nil {
				if hookErr := opts.Hooks.OnPanic(ctx, tx, r); hookErr != nil {
					m.logger.Error("Panic hook failed", "error", hookErr)
				}
			}

			if rollbackErr := tx.RollbackTx(); rollbackErr != nil {
				m.logger.Error("Failed to rollback after panic",
					"tx_id", tx.ID(),
					"error", rollbackErr,
				)
			}

			err = fmt.Errorf("transaction panicked: %v", r)
			atomic.AddInt64(&m.stats.failedTx, 1)
		}
	}()

	// Execute after begin hook
	if opts.Hooks != nil && opts.Hooks.AfterBegin != nil {
		if err := opts.Hooks.AfterBegin(ctx, tx); err != nil {
			tx.RollbackTx()
			return fmt.Errorf("after begin hook failed: %w", err)
		}
	}

	// Execute the main function
	startTime := time.Now()
	err = fn(tx)
	duration := time.Since(startTime)

	// Update statistics
	atomic.AddInt64(&m.stats.totalTx, 1)
	m.updateAverageExecutionTime(duration)

	if err != nil {
		// Execute before rollback hook
		if opts.Hooks != nil && opts.Hooks.BeforeRollback != nil {
			if hookErr := opts.Hooks.BeforeRollback(ctx, tx); hookErr != nil {
				m.logger.Error("Before rollback hook failed", "error", hookErr)
			}
		}

		// Rollback transaction
		if rollbackErr := tx.RollbackTx(); rollbackErr != nil {
			m.logger.Error("Failed to rollback transaction",
				"tx_id", tx.ID(),
				"original_error", err,
				"rollback_error", rollbackErr,
			)
			err = fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rollbackErr)
		}

		// Execute after rollback hook
		if opts.Hooks != nil && opts.Hooks.AfterRollback != nil {
			if hookErr := opts.Hooks.AfterRollback(ctx, tx, err); hookErr != nil {
				m.logger.Error("After rollback hook failed", "error", hookErr)
			}
		}

		atomic.AddInt64(&m.stats.rolledBackTx, 1)
		return err
	}

	// Execute before commit hook
	if opts.Hooks != nil && opts.Hooks.BeforeCommit != nil {
		if err := opts.Hooks.BeforeCommit(ctx, tx); err != nil {
			tx.RollbackTx()
			return fmt.Errorf("before commit hook failed: %w", err)
		}
	}

	// Commit transaction
	if err := tx.CommitTx(); err != nil {
		atomic.AddInt64(&m.stats.failedTx, 1)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Execute after commit hook
	if opts.Hooks != nil && opts.Hooks.AfterCommit != nil {
		if err := opts.Hooks.AfterCommit(ctx, tx); err != nil {
			m.logger.Error("After commit hook failed", "error", err)
			// Don't return error here as transaction is already committed
		}
	}

	atomic.AddInt64(&m.stats.committedTx, 1)
	return nil
}

// executeWithRetry implements intelligent retry logic
func (m *Manager) executeWithRetry(ctx context.Context, opts TxOptions, fn func(tx Transaction) error) error {
	policy := opts.RetryPolicy
	var lastErr error

	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		err := m.executeOnce(ctx, opts, fn)
		if err == nil {
			if attempt > 1 {
				atomic.AddInt64(&m.stats.retryCount, 1)
				m.logger.Info("Transaction succeeded after retry",
					"attempt", attempt,
					"total_attempts", policy.MaxAttempts,
				)
			}
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !m.isRetryableError(err, policy) {
			break
		}

		// Don't retry on the last attempt
		if attempt == policy.MaxAttempts {
			break
		}

		// Calculate delay
		delay := m.calculateRetryDelay(attempt, policy)

		m.logger.Warn("Transaction failed, retrying",
			"attempt", attempt,
			"total_attempts", policy.MaxAttempts,
			"error", err,
			"delay", delay,
		)

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("transaction failed after %d attempts: %w", policy.MaxAttempts, lastErr)
}

// Begin starts a new transaction manually
func (m *Manager) Begin(ctx context.Context) (Transaction, error) {
	return m.BeginWithOptions(ctx, TxOptions{})
}

// BeginWithOptions starts a new transaction with specific options
func (m *Manager) BeginWithOptions(ctx context.Context, opts TxOptions) (Transaction, error) {
	if m.closed {
		return nil, ErrManagerClosed
	}

	opts = m.applyDefaults(opts)

	tx, err := m.provider.Begin(ctx, opts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// ExecuteBatch executes multiple operations in a single transaction
func (m *Manager) ExecuteBatch(ctx context.Context, operations []BatchOperation) error {
	return m.Execute(ctx, func(tx Transaction) error {
		for i, op := range operations {
			if err := op.Operation(tx); err != nil {
				// Handle operation-specific error handling
				if op.OnError != nil {
					if handledErr := op.OnError(err); handledErr != nil {
						return fmt.Errorf("batch operation %d (%s) failed: %w", i, op.Name, handledErr)
					}
					// Error was handled, continue
					continue
				}
				return fmt.Errorf("batch operation %d (%s) failed: %w", i, op.Name, err)
			}
		}
		return nil
	})
}

// Helper methods and utilities
func (m *Manager) applyDefaults(opts TxOptions) TxOptions {
	if opts.Timeout == 0 {
		opts.Timeout = m.config.DefaultTimeout
	}
	if opts.SavepointPrefix == "" {
		opts.SavepointPrefix = "sp"
	}
	return opts
}

func (m *Manager) canStartTransaction() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.activeTxs) < m.config.MaxConcurrentTx
}

func (m *Manager) addActiveTransaction(id string, txCtx *transactionContext) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeTxs[id] = txCtx
}

func (m *Manager) removeActiveTransaction(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.activeTxs, id)
}

func (m *Manager) isRetryableError(err error, policy *RetryPolicy) bool {
	// Implement sophisticated error detection
	errStr := err.Error()
	for _, retryableErr := range policy.RetryableErrors {
		if contains(errStr, retryableErr) {
			return true
		}
	}

	// Common retryable database errors
	retryablePatterns := []string{
		"deadlock",
		"connection reset",
		"timeout",
		"connection refused",
		"serialization failure",
	}

	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}

	return false
}

func (m *Manager) calculateRetryDelay(attempt int, policy *RetryPolicy) time.Duration {
	switch policy.BackoffStrategy {
	case BackoffLinear:
		delay := time.Duration(attempt) * policy.InitialDelay
		if delay > policy.MaxDelay {
			return policy.MaxDelay
		}
		return delay

	case BackoffExponential:
		delay := time.Duration(math.Pow(2, float64(attempt-1))) * policy.InitialDelay
		if delay > policy.MaxDelay {
			return policy.MaxDelay
		}
		return delay

	case BackoffFixed:
		return policy.InitialDelay

	default:
		return policy.InitialDelay
	}
}

func (m *Manager) updateAverageExecutionTime(duration time.Duration) {
	// Thread-safe average calculation using atomic operations
	for {
		oldAvg := atomic.LoadInt64((*int64)(&m.stats.avgExecutionTime))
		oldCount := atomic.LoadInt64(&m.stats.totalTx)

		if oldCount == 0 {
			if atomic.CompareAndSwapInt64((*int64)(&m.stats.avgExecutionTime), oldAvg, int64(duration)) {
				break
			}
			continue
		}

		newAvg := time.Duration((int64(oldAvg)*oldCount + int64(duration)) / (oldCount + 1))
		if atomic.CompareAndSwapInt64((*int64)(&m.stats.avgExecutionTime), oldAvg, int64(newAvg)) {
			break
		}
	}
}

func (m *Manager) cleanupLoop() {
	for {
		select {
		case <-m.shutdownCh:
			return
		case <-m.cleanupTicker.C:
			m.cleanupExpiredTransactions()
		}
	}
}

func (m *Manager) cleanupExpiredTransactions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, txCtx := range m.activeTxs {
		if now.Sub(txCtx.startTime) > txCtx.options.Timeout {
			m.logger.Warn("Force rolling back expired transaction", "tx_id", id)
			if err := txCtx.tx.Rollback(); err != nil {
				m.logger.Error("Failed to rollback expired transaction", "tx_id", id, "error", err)
			}
			delete(m.activeTxs, id)
		}
	}
}

// Stats returns current transaction statistics
func (m *Manager) Stats() TxStats {
	m.mu.RLock()
	activeTx := int64(len(m.activeTxs))
	m.mu.RUnlock()

	return TxStats{
		ActiveTransactions:     activeTx,
		TotalTransactions:      atomic.LoadInt64(&m.stats.totalTx),
		CommittedTransactions:  atomic.LoadInt64(&m.stats.committedTx),
		RolledBackTransactions: atomic.LoadInt64(&m.stats.rolledBackTx),
		FailedTransactions:     atomic.LoadInt64(&m.stats.failedTx),
		AverageExecutionTime:   time.Duration(atomic.LoadInt64((*int64)(&m.stats.avgExecutionTime))),
		DeadlockCount:          atomic.LoadInt64(&m.stats.deadlockCount),
		RetryCount:             atomic.LoadInt64(&m.stats.retryCount),
	}
}

// Close shuts down the transaction manager gracefully
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	m.closed = true
	close(m.shutdownCh)
	m.cleanupTicker.Stop()

	// Wait for active transactions to complete or force close them
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for len(m.activeTxs) > 0 {
		select {
		case <-timeout:
			// Force close remaining transactions
			for id, txCtx := range m.activeTxs {
				m.logger.Warn("Force closing transaction during shutdown", "tx_id", id)
				txCtx.tx.Rollback()
			}
			return nil
		case <-ticker.C:
			// Continue waiting
		}
	}

	return nil
}

// Internal types
type transactionContext struct {
	tx        Transaction
	startTime time.Time
	context   context.Context
	options   TxOptions
}

type stats struct {
	totalTx          int64
	committedTx      int64
	rolledBackTx     int64
	failedTx         int64
	avgExecutionTime time.Duration
	deadlockCount    int64
	retryCount       int64
}

func newStats() *stats {
	return &stats{}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findInString(s, substr))))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
