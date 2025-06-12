package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ClassRouter struct {
	Controller   *controller.ClassController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewClassRouter(class_controller *controller.ClassController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ClassRouter {
	return &ClassRouter{Controller: class_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ClassRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	cc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/classes", cc.GetAllClasses)
	public.GET("/class/:id", cc.GetClassByID)

	protected := rg.Group("")
	protected.Use(i.Authorized.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/class/create", cc.CreateClass)
	protected.PUT("/class/update/:id", cc.UpdateClass)
	protected.DELETE("/class/:id", cc.DeleteClass)
}
