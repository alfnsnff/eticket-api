package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHarborRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewHarborRepository(db)
	hc := &controller.HarborController{
		HarborUsecase: usecase.NewHarborUsecase(hr),
	}
	group.POST("/harbor", hc.CreateHarbor)
	group.GET("/harbors", hc.GetAllHarbors)
	group.GET("/harbor/:id", hc.GetHarborByID)
	group.PUT("/harbor/:id", hc.UpdateHarbor)
	group.DELETE("/harbor/:id", hc.DeleteHarbor)
}
