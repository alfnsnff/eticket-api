package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewSessionRouter(db *gorm.DB, group *gin.RouterGroup) {
	sr := repository.NewSessionRepository()
	tr := repository.NewTicketRepository()
	scr := repository.NewScheduleRepository()
	ar := repository.NewAllocationRepository()
	mr := repository.NewManifestRepository()
	fr := repository.NewFareRepository()
	hc := &controller.SessionController{
		SessionUsecase: usecase.NewSessionUsecase(db, sr, tr, scr, ar, mr, fr),
	}
	group.POST("/session/create", hc.CreateSession)
	group.POST("/session/ticket/lock", hc.LockTicket)
	group.GET("/sessions", hc.GetAllSessions)
	group.GET("/session/:id", hc.GetSessionByID)
	group.GET("/session/uuid/:sessionid", hc.GetSessionBySessionID)
	group.PUT("/session/update/:id", hc.UpdateSession)
	group.DELETE("/session/:id", hc.DeleteSession)
}
