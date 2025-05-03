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

func (bc *BookingController) CreateBooking(ctx *gin.Context) {
	request := new(model.WriteBookingRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	err := bc.BookingUsecase.CreateBooking(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

func (bc *BookingController) GetAllBookings(ctx *gin.Context) {
	datas, err := bc.BookingUsecase.GetAllBookings(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Bookings retrieved successfully", nil))
}

// GetBookingByID retrieves a booking by its ID
func (bc *BookingController) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	data, err := bc.BookingUsecase.GetBookingByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Booking retrieved successfully", nil))
}

func (bc *BookingController) UpdateBooking(ctx *gin.Context) {
	request := new(model.UpdateBookingRequest)
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Booking ID is required", nil))
		return
	}

	err := bc.BookingUsecase.UpdateBooking(ctx, uint(id), request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

func (bc *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := bc.BookingUsecase.DeleteBooking(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}
