package middleware

import (
	repository "eticket-api/internal/repository/auth"
	"eticket-api/pkg/jwt"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenService     *jwt.TokenManager
	UserRepository   *repository.UserRepository
	RefreshTokenRepo *repository.AuthRepository
}

func NewAuthMiddleware(token_manager *jwt.TokenManager, user_repository *repository.UserRepository, auth_repository *repository.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{
		TokenService:     token_manager,
		UserRepository:   user_repository,
		RefreshTokenRepo: auth_repository,
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
		claims, err := am.TokenService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.NewErrorResponse("Invalid or expired token", err.Error()))
			return
		}

		// Inject user info into Gin context
		c.Set("userID", claims.ID)
		c.Set("username", claims.Username)
		c.Set("token", tokenStr)

		c.Next()
	}
}
