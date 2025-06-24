package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"eticket-api/internal/common/logger"

	"github.com/gin-gonic/gin"
)

type RecoveryMiddleware struct {
	Log logger.Logger
}

func NewRecoveryMiddleware(log logger.Logger) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		Log: log,
	}
}

func (rm *RecoveryMiddleware) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log panic with stack trace
				rm.Log.WithFields(map[string]interface{}{
					"panic":   fmt.Sprintf("%v", rec),
					"stack":   string(debug.Stack()),
					"path":    c.Request.URL.Path,
					"method":  c.Request.Method,
					"headers": c.Request.Header,
				}).Error("Unhandled panic occurred")

				// Respond with generic 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}

// Legacy function for backward compatibility
func Recovery(log logger.Logger) gin.HandlerFunc {
	middleware := NewRecoveryMiddleware(log)
	return middleware.Set()
}
