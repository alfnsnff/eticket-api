package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/injector"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewSessionRouter(ic *injector.Container, rg *gin.RouterGroup) {
	csr := ic.Repository.SessionRepository
	tr := ic.Repository.TicketRepository
	scr := ic.Repository.ScheduleRepository
	ar := ic.Repository.AllocationRepository
	mr := ic.Repository.ManifestRepository
	fr := ic.Repository.FareRepository
	br := ic.Repository.BookingRepository
	sc := &controller.SessionController{
		SessionUsecase: usecase.NewSessionUsecase(ic.Tx, csr, tr, scr, ar, mr, fr, br),
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
