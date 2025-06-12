package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"eticket-api/internal/common/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type AuthorizeMiddleware struct {
	Enforcer *casbin.Enforcer
}

func NewAuthorizeMiddleware(enforcer *casbin.Enforcer) *AuthorizeMiddleware {
	return &AuthorizeMiddleware{Enforcer: enforcer}
}

// Gin middleware to enforce RBAC
func (i *AuthorizeMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user role from context (set during auth middleware)
		role, exists := c.Get("rolename")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Role not found in context", nil))
			return
		}

		obj := strings.TrimPrefix(c.FullPath(), "/api")
		act := strings.ToUpper(c.Request.Method)

		allowed, err := i.Enforcer.Enforce(role, obj, act)
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
