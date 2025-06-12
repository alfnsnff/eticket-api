package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type AllocationRouter struct {
	Controller   *controller.AllocationController
	Authenticate *middleware.AuthenticateMiddleware
	Authorized   *middleware.AuthorizeMiddleware
}

func NewAllocationRouter(allocation_controller *controller.AllocationController, authtenticate *middleware.AuthenticateMiddleware, authorized *middleware.AuthorizeMiddleware) *AllocationRouter {
	return &AllocationRouter{Controller: allocation_controller, Authenticate: authtenticate, Authorized: authorized}
}

func (i AllocationRouter) Set(router *gin.Engine, rg *gin.RouterGroup) {

	public := rg.Group("") // No middleware
	public.GET("/allocations", i.Controller.GetAllAllocations)
	public.GET("/allocation/:id", i.Controller.GetAllocationByID)

	protected := rg.Group("")
	protected.Use(i.Authenticate.Handle())
	// protected.Use(i.Authorized.Handle())

	protected.POST("/allocation/create", i.Controller.CreateAllocation)
	protected.PUT("/allocation/update/:id", i.Controller.UpdateAllocation)
	protected.DELETE("/allocation/:id", i.Controller.DeleteAllocation)
}
