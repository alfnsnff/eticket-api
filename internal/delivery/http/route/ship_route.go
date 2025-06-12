package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type ShipRouter struct {
	Controller   *controller.ShipController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewShipRouter(ship_controller *controller.ShipController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *ShipRouter {
	return &ShipRouter{Controller: ship_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i ShipRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	shc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/ships", shc.GetAllShips)
	public.GET("/ships/:id", shc.GetShipByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)
}
