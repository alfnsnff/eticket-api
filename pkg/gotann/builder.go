package gotann

import (
	"log/slog"
	"time"
)

// Builder provides a fluent interface for creating transaction managers
type Builder struct {
	provider Provider
	config   Config
	logger   *slog.Logger
}

// NewBuilder creates a new transaction manager builder
func NewBuilder(provider Provider) *Builder {
	return &Builder{
		provider: provider,
		config:   DefaultConfig(),
	}
}

// WithMaxConcurrentTransactions sets the maximum concurrent transactions
func (b *Builder) WithMaxConcurrentTransactions(max int) *Builder {
	b.config.MaxConcurrentTx = max
	return b
}

// WithDefaultTimeout sets the default transaction timeout
func (b *Builder) WithDefaultTimeout(timeout time.Duration) *Builder {
	b.config.DefaultTimeout = timeout
	return b
}

// WithLogger sets the logger
func (b *Builder) WithLogger(logger *slog.Logger) *Builder {
	b.logger = logger
	return b
}

// WithMetrics enables or disables metrics collection
func (b *Builder) WithMetrics(enabled bool) *Builder {
	b.config.EnableMetrics = enabled
	return b
}

// WithTracing enables or disables distributed tracing
func (b *Builder) WithTracing(enabled bool) *Builder {
	b.config.EnableTracing = enabled
	return b
}

// Build creates the transaction manager
func (b *Builder) Build() *Manager {
	return NewManager(b.provider, b.config, b.logger)
}

// QuickOptions provides quick configuration for common scenarios
type QuickOptions struct {
	Timeout  time.Duration
	Retries  int
	ReadOnly bool
}

// Quick creates a transaction manager with minimal configuration
func Quick(provider Provider, opts QuickOptions) *Manager {
	config := DefaultConfig()
	if opts.Timeout > 0 {
		config.DefaultTimeout = opts.Timeout
	}
	if opts.Retries > 0 {
		config.MaxRetryAttempts = opts.Retries
	}

	return NewManager(provider, config, nil)
}

// WithRetry creates TxOptions with retry policy
func WithRetry(maxAttempts int, initialDelay time.Duration) TxOptions {
	return TxOptions{
		RetryPolicy: &RetryPolicy{
			MaxAttempts:     maxAttempts,
			InitialDelay:    initialDelay,
			MaxDelay:        30 * time.Second,
			BackoffStrategy: BackoffExponential,
			RetryableErrors: []string{"deadlock", "timeout", "connection"},
		},
	}
}

// WithTimeout creates TxOptions with timeout
func WithTimeout(timeout time.Duration) TxOptions {
	return TxOptions{
		Timeout: timeout,
	}
}

// WithIsolation creates TxOptions with isolation level
func WithIsolation(level IsolationLevel) TxOptions {
	return TxOptions{
		IsolationLevel: level,
	}
}

// ReadOnly creates TxOptions for read-only transactions
func ReadOnly() TxOptions {
	return TxOptions{
		ReadOnly: true,
	}
}

// Combine merges multiple TxOptions
func Combine(opts ...TxOptions) TxOptions {
	result := TxOptions{}

	for _, opt := range opts {
		if opt.Timeout > 0 {
			result.Timeout = opt.Timeout
		}
		if opt.IsolationLevel != IsolationDefault {
			result.IsolationLevel = opt.IsolationLevel
		}
		if opt.ReadOnly {
			result.ReadOnly = true
		}
		if opt.RetryPolicy != nil {
			result.RetryPolicy = opt.RetryPolicy
		}
		if opt.Hooks != nil {
			result.Hooks = opt.Hooks
		}
		if opt.Context != nil {
			if result.Context == nil {
				result.Context = make(map[string]interface{})
			}
			for k, v := range opt.Context {
				result.Context[k] = v
			}
		}
		if opt.SavepointPrefix != "" {
			result.SavepointPrefix = opt.SavepointPrefix
		}
		if opt.DisableLogging {
			result.DisableLogging = true
		}
	}

	return result
}
