package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	authrepository "eticket-api/internal/repository/auth"
	authusecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRoleRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := authrepository.NewRoleRepository()
	hc := &authcontroller.RoleController{
		RoleUsecase: authusecase.NewRoleUsecase(db, rr),
	}
	group.POST("/role/create", hc.CreateRole)
}
