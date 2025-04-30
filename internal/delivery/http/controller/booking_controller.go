package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	BookingUsecase *usecase.BookingUsecase
}

// NewBookingController creates a new BookingController instance
func NewBookingController(booking_usecase *usecase.BookingUsecase) *BookingController {
	return &BookingController{BookingUsecase: booking_usecase}
}

// CreateBooking handles creating a new booking
func (h *BookingController) CreateBooking(ctx *gin.Context) {
	request := new(model.WriteBookingRequest)
	// Bind request body to DTO
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	// booking, _ := mapper.ToBookingEntity(request)

	// Call use case
	err := h.BookingUsecase.CreateBooking(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

// // CreateBookingWithTickets handles creating a booking with associated tickets
// func (h *BookingController) CreateBookingWithTickets(ctx *gin.Context) {
// 	var bookingCreate dto.BookingCreate

// 	if err := ctx.ShouldBindJSON(&bookingCreate); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	booking, tickets := dto.ToBookingEntity(&bookingCreate)

// 	// Call use case
// 	err := h.BookingUsecase.CreateBookingWithTickets(ctx, &booking, &tickets)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking with tickets", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking with tickets created successfully", nil))
// }

// GetAllBookings retrieves all bookings
func (h *BookingController) GetAllBookings(ctx *gin.Context) {
	datas, err := h.BookingUsecase.GetAllBookings(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	// bookingDTOs := dto.ToBookingDTOs(bookings)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Bookings retrieved successfully", nil))
}

// GetBookingByID retrieves a booking by its ID
func (h *BookingController) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	data, err := h.BookingUsecase.GetBookingByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	// bookingDTO := dto.ToBookingDTO(booking)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Booking retrieved successfully", nil))
}

// UpdateBooking updates an existing booking
func (h *BookingController) UpdateBooking(ctx *gin.Context) {
	request := new(model.WriteBookingRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Booking ID is required", nil))
		return
	}

	// booking, _ := dto.ToBookingEntity(request)

	err := h.BookingUsecase.UpdateBooking(ctx, uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

// DeleteBooking deletes a booking by its ID
func (h *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := h.BookingUsecase.DeleteBooking(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}
