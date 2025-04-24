package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewShipClassRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewShipClassRepository(db)
	hc := &controller.ShipClassController{
		ShipClassUsecase: usecase.NewShipClassUsecase(hr),
	}
	group.POST("/shipClass", hc.CreateShipClass)
	group.GET("/shipClasses", hc.GetAllShipClasses)
	group.GET("/shipClasses/ship/:id", hc.GetShipClassByShipID)
	group.GET("/shipClass/:id", hc.GetShipClassByID)
	group.PUT("/shipClass/:id", hc.UpdateShipClass)
	group.DELETE("/shipClass/:id", hc.DeleteShipClass)
}
