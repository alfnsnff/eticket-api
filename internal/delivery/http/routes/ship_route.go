package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ShipRoute struct {
	Controller   *controller.ShipController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewShipRoute(ship_controller *controller.ShipController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ShipRoute {
	return &ShipRoute{Controller: ship_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ShipRoute) Set(router *gin.Engine) {
	shc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/ships", shc.GetAllShips)
	public.GET("/ships/:id", shc.GetShipByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)
}
