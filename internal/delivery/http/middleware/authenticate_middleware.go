package middleware

import (
	"eticket-api/internal/common/jwt"
	"eticket-api/internal/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticateMiddleware struct {
	TokenUtil *jwt.TokenUtil
}

func NewAuthenticateMiddleware(token_util *jwt.TokenUtil) *AuthenticateMiddleware {
	return &AuthenticateMiddleware{
		TokenUtil: token_util,
	}
}

// Middleware method to authenticate access token via cookie
func (am *AuthenticateMiddleware) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from cookie
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Missing token", err.Error()))
			return
		}

		// Validate token
		claims, err := am.TokenUtil.ValidateToken(tokenStr)
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
