package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type FareRouter struct {
	Controller   *controller.FareController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewFareRouter(fare_controller *controller.FareController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *FareRouter {
	return &FareRouter{Controller: fare_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i FareRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	fc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/fares", fc.GetAllFares)
	public.GET("/fare/:id", fc.GetFareByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/fare/create", fc.CreateFare)
	protected.PUT("/fare/update/:id", fc.UpdateFare)
	protected.DELETE("/fare/:id", fc.DeleteFare)
}
