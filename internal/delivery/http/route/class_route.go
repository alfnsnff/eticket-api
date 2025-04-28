package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewClassRouter(db *gorm.DB, group *gin.RouterGroup) {
	cr := repository.NewClassRepository()
	cc := &controller.ClassController{
		ClassUsecase: usecase.NewClassUsecase(db, cr),
	}
	group.POST("/class", cc.CreateClass)
	group.GET("/classes", cc.GetAllClasses)
	group.GET("/class/:id", cc.GetClassByID)
	group.PUT("/class/:id", cc.UpdateClass)
	group.DELETE("/class/:id", cc.DeleteClass)
}
