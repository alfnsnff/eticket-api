package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewClassRouter(db *gorm.DB, group *gin.RouterGroup) {
	cr := repository.NewClassRepository(db)
	cc := &controller.ClassController{
		ClassUsecase: usecase.NewClassUsecase(cr),
	}
	group.POST("/classes", cc.CreateClass)
	group.GET("/classes", cc.GetAllClasses)
	group.GET("/classes/:id", cc.GetClassByID)
	group.PUT("/classes/:id", cc.UpdateClass)
	group.DELETE("/classes/:id", cc.DeleteClass)
}
