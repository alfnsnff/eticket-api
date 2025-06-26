package middleware

import (
	"eticket-api/internal/common/token"
	"eticket-api/internal/delivery/http/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware method to authenticate access token via cookie
func Authenticate(token_util token.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Missing token", err.Error()))
			return
		}

		// Validate token
		claims, err := token_util.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid or expired token", err.Error()))
			return
		}

		c.Set("rolename", claims.User.Role.RoleName)
		c.Set("token", tokenStr)

		c.Next()
	}
}
