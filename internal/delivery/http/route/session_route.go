package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewSessionRouter(ic *injector.Container, rg *gin.RouterGroup) {
	csr := ic.SessionRepository
	tr := ic.TicketRepository
	scr := ic.ScheduleRepository
	ar := ic.AllocationRepository
	mr := ic.ManifestRepository
	fr := ic.FareRepository
	sc := &controller.SessionController{
		SessionUsecase: usecase.NewSessionUsecase(ic.SessionUsecase.Tx, csr, tr, scr, ar, mr, fr),
	}

	public := rg.Group("") // No middleware
	public.POST("/session/ticket/lock", sc.SessionTicketLock)
	public.POST("/session/ticket/data/entry", sc.SessionTicketDataEntry)
	public.GET("/sessions", sc.GetAllSessions)
	public.GET("/session/:id", sc.GetSessionByID)
	public.GET("/session/uuid/:sessionid", sc.GetSessionBySessionID)

	protected := rg.Group("")
	middleware := middleware.NewAuthMiddleware(ic.TokenManager, ic.UserRepository, ic.AuthRepository)
	protected.Use(middleware.Authenticate())

	protected.POST("/session/create", sc.CreateSession)
	protected.PUT("/session/update/:id", sc.UpdateSession)
	protected.DELETE("/session/:id", sc.DeleteSession)
}
