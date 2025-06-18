package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	Controller   *controller.AuthController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewAuthRoute(auth_controller *controller.AuthController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *AuthRoute {
	return &AuthRoute{Controller: auth_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i AuthRoute) Set(router *gin.Engine) {

	public := router.Group("") // No middleware
	public.GET("/auth/me", i.Controller.Me)
	public.POST("/auth/login", i.Controller.Login)
	public.POST("/auth/refresh", i.Controller.RefreshToken)
	public.POST("/auth/forget-password", i.Controller.ForgetPassword)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/auth/logout", i.Controller.Logout)
}
