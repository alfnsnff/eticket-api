package errors

import (
	"strings"
)

func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	// GORM wraps the driver error, so check the error string
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "Duplicate entry")
}
