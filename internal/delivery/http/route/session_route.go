package route

import (
	"eticket-api/internal/delivery/http/controller"
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
	public.POST("/session/create", sc.CreateSession)
	public.PUT("/session/update/:id", sc.UpdateSession)
	public.DELETE("/session/:id", sc.DeleteSession)
}
