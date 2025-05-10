package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewHarborRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewHarborRepository()
	hc := &controller.HarborController{
		HarborUsecase: usecase.NewHarborUsecase(db, hr),
	}
	group.POST("/harbor/create", hc.CreateHarbor)
	group.GET("/harbors", hc.GetAllHarbors)
	group.GET("/harbor/:id", hc.GetHarborByID)
	group.PUT("/harbor/update/:id", hc.UpdateHarbor)
	group.DELETE("/harbor/:id", hc.DeleteHarbor)
}
