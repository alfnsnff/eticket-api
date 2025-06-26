package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"eticket-api/internal/common/enforcer"
	"eticket-api/internal/delivery/http/response"

	"github.com/gin-gonic/gin"
)

// Gin middleware to enforce RBAC
func Authorize(enforcer enforcer.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user role from context (set during auth middleware)
		role, exists := c.Get("rolename")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Role not found in context", nil))
			return
		}

		obj := strings.TrimPrefix(c.FullPath(), "/api")
		act := strings.ToUpper(c.Request.Method)

		allowed, err := enforcer.Enforce(fmt.Sprint(role), obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewErrorResponse("Authorization error", nil))
			return
		}
		log.Printf("Checking permissions for role: %v  : %v", role, allowed)
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, response.NewErrorResponse(
				fmt.Sprintf("Role %v does not have %v permission for %v", role, act, obj), nil))
			return
		}

		c.Next()
	}
}
