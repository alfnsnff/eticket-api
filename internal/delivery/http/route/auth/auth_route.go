package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	usecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(ic *injector.Container, rg *gin.RouterGroup) {
	ar := ic.AuthRepository
	ur := ic.UserRepository
	tm := ic.TokenManager
	auc := &authcontroller.AuthController{
		Cfg:          ic.Cfg,
		TokenManager: ic.TokenManager,
		AuthUsecase:  usecase.NewAuthUsecase(ic.AuthUsecase.Tx, ar, ur, tm),
	}

	public := rg.Group("") // No middleware
	public.POST("/auth/login", auc.Login)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/auth/logout", auc.Logout)
}
