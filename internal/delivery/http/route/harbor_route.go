package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type HarborRouter struct {
	Controller   *controller.HarborController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewHarborRouter(harbor_controller *controller.HarborController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *HarborRouter {
	return &HarborRouter{Controller: harbor_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i HarborRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {
	hc := i.Controller

	public := rg.Group("") // No middleware
	public.GET("/harbors", hc.GetAllHarbors)
	public.GET("/harbor/:id", hc.GetHarborByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/harbor/create", hc.CreateHarbor)
	protected.PUT("/harbor/update/:id", hc.UpdateHarbor)
	protected.DELETE("/harbor/:id", hc.DeleteHarbor)
}
