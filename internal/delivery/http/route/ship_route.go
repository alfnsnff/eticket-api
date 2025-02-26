package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewShipRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewShipRepository(db)
	hc := &controller.ShipController{
		ShipUsecase: usecase.NewShipUsecase(hr),
	}
	group.POST("/ship", hc.CreateShip)
	group.GET("/ships", hc.GetAllShips)

	group.GET("/ship/:id", hc.GetShipByID)
	group.PUT("/ship/:id", hc.UpdateShip)
	group.DELETE("/ship/:id", hc.DeleteShip)
}
