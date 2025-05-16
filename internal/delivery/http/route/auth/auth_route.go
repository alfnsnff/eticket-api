package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	authrepository "eticket-api/internal/repository/auth"
	authusecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthRouter(db *gorm.DB, public *gin.RouterGroup, protect *gin.RouterGroup) {
	ar := authrepository.NewAuthRepository()
	ur := authrepository.NewUserRepository()
	hc := &authcontroller.AuthController{
		AuthUsecase: authusecase.NewAuthUsecase(db, ar, ur),
	}
	public.POST("/auth/login", hc.Login)
	protect.POST("/auth/logout", hc.Logout)
}
