package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewClassRouter(ic *injector.Container, rg *gin.RouterGroup) {
	cr := ic.Repository.ClassRepository
	cc := &controller.ClassController{
		ClassUsecase: usecase.NewClassUsecase(ic.Tx, cr),
	}

	public := rg.Group("") // No middleware
	public.GET("/classes", cc.GetAllClasses)
	public.GET("/class/:id", cc.GetClassByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager)
	protected.Use(middleware.Authenticate())

	protected.POST("/class/create", cc.CreateClass)
	protected.PUT("/class/update/:id", cc.UpdateClass)
	protected.DELETE("/class/:id", cc.DeleteClass)
}
