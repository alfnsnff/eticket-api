package middleware

import (
	"eticket-api/pkg/utils/helper/auth"
	"eticket-api/pkg/utils/helper/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authenticate validates the JWT token and injects user info into context.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Missing Authorization header", nil))
			return
		}

		// More efficient split for "Bearer <token>"
		scheme, token, found := strings.Cut(authHeader, " ")
		if !found || strings.ToLower(scheme) != "bearer" || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Authorization header must be in 'Bearer <token>' format", nil))
			return
		}

		// Validate JWT token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid or expired token", err.Error()))
			return
		}

		// Inject values into Gin context
		c.Set("userID", claims.ID)
		c.Set("username", claims.Username)
		c.Set("token", token) // Optional: allow downstream logic access to raw token

		c.Next()
	}
}
