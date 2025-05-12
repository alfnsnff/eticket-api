package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAllocationRouter(db *gorm.DB, group *gin.RouterGroup) {
	ar := repository.NewAllocationRepository()
	sr := repository.NewScheduleRepository()
	fr := repository.NewFareRepository()
	hc := &controller.AllocationController{
		AllocationUsecase: usecase.NewAllocationUsecase(db, ar, sr, fr),
	}
	group.POST("/allocation/create", hc.CreateAllocation)
	group.GET("/allocations", hc.GetAllAllocations)
	group.GET("/allocation/:id", hc.GetAllocationByID)
	group.PUT("/allocation/update/:id", hc.UpdateAllocation)
	group.DELETE("/allocation/:id", hc.DeleteAllocation)
}
