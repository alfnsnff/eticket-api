package middleware

import (
	"eticket-api/internal/common/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggerMiddleware struct {
	Log logger.Logger
}

func NewLoggerMiddleware(log logger.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{
		Log: log,
	}
}

func (lm *LoggerMiddleware) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path += "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		lm.Log.WithFields(map[string]interface{}{
			"status":     status,
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"latency":    latency.String(),
		}).Info("incoming request")
	}
}

// Legacy function for backward compatibility
func Logger(log logger.Logger) gin.HandlerFunc {
	middleware := NewLoggerMiddleware(log)
	return middleware.Set()
}
