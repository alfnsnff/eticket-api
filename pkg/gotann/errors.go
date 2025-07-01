package gotann

import (
	"errors"
	"strings"
)

var (
	// Core errors
	ErrManagerClosed        = errors.New("transaction manager is closed")
	ErrTooManyTransactions  = errors.New("too many concurrent transactions")
	ErrTransactionNotActive = errors.New("transaction is not active")
	ErrSavepointNotFound    = errors.New("savepoint not found")
	ErrUnsupportedFeature   = errors.New("feature not supported by provider")

	// Retry errors
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded")
	ErrRetryTimeout       = errors.New("retry timeout exceeded")

	// Configuration errors
	ErrInvalidConfiguration = errors.New("invalid configuration")
	ErrProviderNotSupported = errors.New("provider not supported")
)

// IsRetryableError checks if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	retryablePatterns := []string{
		"deadlock",
		"connection reset",
		"timeout",
		"connection refused",
		"serialization failure",
		"could not serialize",
		"retry transaction",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}
