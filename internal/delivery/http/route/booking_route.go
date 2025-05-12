package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewBookingRouter(db *gorm.DB, group *gin.RouterGroup) {
	br := repository.NewBookingRepository()
	tr := repository.NewTicketRepository()
	sr := repository.NewSessionRepository()
	hc := &controller.BookingController{
		BookingUsecase: usecase.NewBookingUsecase(db, br, tr, sr),
	}
	group.POST("/booking/create", hc.CreateBooking)
	group.POST("/booking/confirm", hc.ConfirmBooking)
	group.GET("/bookings", hc.GetAllBookings)
	group.GET("/booking/:id", hc.GetBookingByID)
	group.PUT("/booking/update/:id", hc.UpdateBooking)
	group.DELETE("/booking/:id", hc.DeleteBooking)
}
