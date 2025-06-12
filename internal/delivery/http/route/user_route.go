package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	Controller   *controller.UserController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewUserRouter(user_controller *controller.UserController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *UserRouter {
	return &UserRouter{Controller: user_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i UserRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	uc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/users", uc.GetAllUsers)
	public.GET("/user/:id", uc.GetUserByID)
	public.POST("/user/create", uc.CreateUser)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.PUT("/user/update/:id", uc.UpdateUser)
	protected.DELETE("/user/:id", uc.DeleteUser)
}
