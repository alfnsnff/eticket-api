package casbin

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func Middleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userID")
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := e.Enforce(userID, obj, act)
		if err != nil || !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
