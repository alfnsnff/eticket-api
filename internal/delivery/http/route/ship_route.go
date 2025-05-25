package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewShipRouter(ic *injector.Container, rg *gin.RouterGroup) {
	repo := repository.NewShipRepository()
	shc := &controller.ShipController{ShipUsecase: usecase.NewShipUsecase(ic.Tx, repo)}

	public := rg.Group("") // No middleware
	public.GET("/ships", shc.GetAllShips)
	public.GET("/ships/:id", shc.GetShipByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager)
	protected.Use(middleware.Authenticate())

	protected.POST("/ship/create", shc.CreateShip)
	protected.PUT("/ship/:id", shc.UpdateShip)
	protected.DELETE("/ship/:id", shc.DeleteShip)
}
