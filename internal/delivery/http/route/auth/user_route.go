package route

import (
	authcontroller "eticket-api/internal/delivery/http/controller/auth"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	usecase "eticket-api/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(ic *injector.Container, rg *gin.RouterGroup) {
	ur := ic.UserRepository
	uc := &authcontroller.UserController{
		UserUsecase: usecase.NewUserUsecase(ic.UserUsecase.Tx, ur),
	}

	public := rg.Group("") // No middleware
	public.GET("/users", uc.GetAllUsers)
	public.GET("/user/:id", uc.GetUserByID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/user/create", uc.CreateUser)
	protected.PUT("/user/update/:id", uc.UpdateUser)
	protected.DELETE("/user/:id", uc.DeleteUser)
}
