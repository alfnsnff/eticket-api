package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewShipRouter(db *gorm.DB, public *gin.RouterGroup, protect *gin.RouterGroup) {
	hr := repository.NewShipRepository()
	hc := &controller.ShipController{
		ShipUsecase: usecase.NewShipUsecase(db, hr),
	}
	protect.POST("/ship/create", hc.CreateShip)
	protect.GET("/ships", hc.GetAllShips)
	protect.GET("/ship/:id", hc.GetShipByID)
	protect.PUT("/ship/update/:id", hc.UpdateShip)
	protect.DELETE("/ship/:id", hc.DeleteShip)
}
