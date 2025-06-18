package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	Controller   *controller.UserController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewUserRoute(user_controller *controller.UserController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *UserRoute {
	return &UserRoute{Controller: user_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i UserRoute) Set(router *gin.Engine) {
	uc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/users", uc.GetAllUsers)
	public.GET("/user/:id", uc.GetUserByID)
	public.POST("/user/create", uc.CreateUser)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.PUT("/user/update/:id", uc.UpdateUser)
	protected.DELETE("/user/:id", uc.DeleteUser)
}
