package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type FareRoute struct {
	Controller   *controller.FareController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewFareRoute(fare_controller *controller.FareController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *FareRoute {
	return &FareRoute{Controller: fare_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i FareRoute) Set(router *gin.Engine) {
	fc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/fares", fc.GetAllFares)
	public.GET("/fare/:id", fc.GetFareByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/fare/create", fc.CreateFare)
	protected.PUT("/fare/update/:id", fc.UpdateFare)
	protected.DELETE("/fare/:id", fc.DeleteFare)
}
