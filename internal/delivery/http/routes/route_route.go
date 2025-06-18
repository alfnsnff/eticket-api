package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RouteRoute struct {
	Controller   *controller.RouteController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewRouteRoute(route_controller *controller.RouteController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *RouteRoute {
	return &RouteRoute{Controller: route_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i RouteRoute) Set(router *gin.Engine) {
	rc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/routes", rc.GetAllRoutes)
	public.GET("/route/:id", rc.GetRouteByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/route/create", rc.CreateRoute)
	protected.PUT("/route//update:id", rc.UpdateRoute)
	protected.DELETE("/route/:id", rc.DeleteRoute)
}
