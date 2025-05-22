package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouteRouter(ic *injector.Container, rg *gin.RouterGroup) {
	rr := ic.Repository.RouteRepository
	rc := &controller.RouteController{
		RouteUsecase: usecase.NewRouteUsecase(ic.Tx, rr),
	}

	public := rg.Group("") // No middleware
	public.GET("/routes", rc.GetAllRoutes)
	public.GET("/route/:id", rc.GetRouteByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.Repository.UserRepository, ic.Repository.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/route/create", rc.CreateRoute)
	protected.PUT("/route//update:id", rc.UpdateRoute)
	protected.DELETE("/route/:id", rc.DeleteRoute)
}
