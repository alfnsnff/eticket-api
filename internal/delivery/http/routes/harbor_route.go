package routes

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type HarborRoute struct {
	Controller   *controller.HarborController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewHarborRoute(harbor_controller *controller.HarborController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *HarborRoute {
	return &HarborRoute{Controller: harbor_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i HarborRoute) Set(router *gin.Engine) {
	hc := i.Controller

	public := router.Group("") // No middleware
	public.GET("/harbors", hc.GetAllHarbors)
	public.GET("/harbor/:id", hc.GetHarborByID)

	protected := router.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)
}
