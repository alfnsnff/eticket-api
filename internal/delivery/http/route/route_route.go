package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RouteRouter struct {
	Controller   *controller.RouteController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewRouteRouter(route_controller *controller.RouteController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *RouteRouter {
	return &RouteRouter{Controller: route_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i RouteRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	rc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/routes", rc.GetAllRoutes)
	public.GET("/route/:id", rc.GetRouteByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/route/create", rc.CreateRoute)
	protected.PUT("/route//update:id", rc.UpdateRoute)
	protected.DELETE("/route/:id", rc.DeleteRoute)
}
