package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewPrebookRouter(db *gorm.DB, group *gin.RouterGroup) {
	ar := repository.NewAllocationRepository()
	sr := repository.NewScheduleRepository()
	fr := repository.NewFareRepository()
	mr := repository.NewManifestRepository()
	tr := repository.NewTicketRepository()
	pc := &controller.PrebookController{
		PreBookUsecase: usecase.NewPrebookTicketsUsecase(db, sr, ar, tr, mr, fr),
	}
	group.POST("/prebook", pc.LockTicket)

}
