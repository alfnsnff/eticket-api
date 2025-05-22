package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewAllocationRouter(ic *injector.Container, rg *gin.RouterGroup) {
	ar := ic.Repository.AllocationRepository
	scr := ic.Repository.ScheduleRepository
	fr := ic.Repository.FareRepository
	ac := &controller.AllocationController{
		AllocationUsecase: usecase.NewAllocationUsecase(ic.Tx, ar, scr, fr),
	}

	public := rg.Group("") // No middleware
	public.GET("/allocations", ac.GetAllAllocations)
	public.GET("/allocation/:id", ac.GetAllocationByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.Repository.UserRepository, ic.Repository.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/allocation/create", ac.CreateAllocation)
	protected.PUT("/allocation/update/:id", ac.UpdateAllocation)
	protected.DELETE("/allocation/:id", ac.DeleteAllocation)
}
