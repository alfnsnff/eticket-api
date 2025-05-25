package middleware

import (
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenManager *jwt.TokenManager
}

func NewAuthMiddleware(token_manager *jwt.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		TokenManager: token_manager,
	}
}

// Middleware method to authenticate access token via cookie
func (am *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Missing token", err.Error()))
			return
		}

		// Validate token
		claims, err := am.TokenManager.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid or expired token", err.Error()))
			return
		}

		// Inject user info into Gin context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roleID", claims.RoleID)
		c.Set("rolename", claims.Rolename)
		c.Set("token", tokenStr)

		c.Next()
	}
}
