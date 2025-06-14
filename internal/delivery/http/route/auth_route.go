package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	Controller   *controller.AuthController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewAuthRouter(auth_controller *controller.AuthController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *AuthRouter {
	return &AuthRouter{Controller: auth_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i AuthRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {

	public := rg.Group("") // No middleware
	public.GET("/auth/me", i.Controller.Me)
	public.POST("/auth/login", i.Controller.Login)
	public.POST("/auth/refresh", i.Controller.RefreshToken)
	public.POST("/auth/forget-password", i.Controller.ForgetPassword)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/auth/logout", i.Controller.Logout)
}
