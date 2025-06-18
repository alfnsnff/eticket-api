package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ClassRoute struct {
	Controller   *controller.ClassController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewClassRoute(class_controller *controller.ClassController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ClassRoute {
	return &ClassRoute{Controller: class_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ClassRoute) Set(router *gin.Engine) {
	cc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/classes", cc.GetAllClasses)
	public.GET("/class/:id", cc.GetClassByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/class/create", cc.CreateClass)
	protected.PUT("/class/update/:id", cc.UpdateClass)
	protected.DELETE("/class/:id", cc.DeleteClass)
}
