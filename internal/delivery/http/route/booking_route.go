package route

import (
	"eticket-api/internal/delivery/http/controller"
	"eticket-api/internal/repository"
	"eticket-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewBookingRouter(db *gorm.DB, group *gin.RouterGroup) {
	hr := repository.NewBookingRepository()
	tr := repository.NewTicketRepository()
	hc := &controller.BookingController{
		BookingUsecase: usecase.NewBookingUsecase(db, hr, tr),
	}
	group.POST("/booking", hc.CreateBooking)
	// group.POST("/booking/with-tickets", hc.CreateBookingWithTickets)
	group.GET("/bookings", hc.GetAllBookings)
	group.GET("/booking/:id", hc.GetBookingByID)
	group.PUT("/booking/:id", hc.UpdateBooking)
	group.DELETE("/booking/:id", hc.DeleteBooking)
}
