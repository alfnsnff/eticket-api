package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	authrepository "eticket-api/internal/repository/auth"
	authusecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserRoleRouter(db *gorm.DB, group *gin.RouterGroup) {
	rr := authrepository.NewRoleRepository()
	ur := authrepository.NewUserRepository()
	usr := authrepository.NewUserRoleRepository()
	hc := &authcontroller.UserRoleController{
		UserRoleUsecase: authusecase.NewUserRoleUsecase(db, rr, ur, usr),
	}
	group.POST("/user/role/assign", hc.CreateUserRole)
	group.GET("/user/roles", hc.GetAllUserRoles)
	group.GET("/user/role/:id", hc.GetUserRoleByID)
	group.PUT("/user/role/update/:id", hc.UpdateUserRole)
	group.DELETE("/user/role/:id", hc.DeleteUserRole)
}
