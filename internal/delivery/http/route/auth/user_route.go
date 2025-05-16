package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	authrepository "eticket-api/internal/repository/auth"
	authusecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserRouter(db *gorm.DB, group *gin.RouterGroup) {
	ur := authrepository.NewUserRepository()
	hc := &authcontroller.UserController{
		UserUsecase: authusecase.NewUserUsecase(db, ur),
	}
	group.POST("/user/create", hc.CreateUser)
	group.POST("/user/login", hc.Login)
}
