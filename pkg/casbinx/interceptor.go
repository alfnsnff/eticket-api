package casbinx

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"eticket-api/pkg/utils/helper/response"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

type Interceptor struct {
	Enforcer *casbin.Enforcer
}

func NewInterceptor(enforcer *casbin.Enforcer) *Interceptor {
	return &Interceptor{Enforcer: enforcer}
}

// Gin middleware to enforce RBAC
func (i *Interceptor) Authorize() gin.HandlerFunc {
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
