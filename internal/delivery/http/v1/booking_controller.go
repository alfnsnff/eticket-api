package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	Validate       validator.Validator
	Log            logger.Logger
	BookingUsecase *usecase.BookingUsecase
}

// NewBookingController creates a new BookingController instance
func NewBookingController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	booking_usecase *usecase.BookingUsecase,

) {
	bc := &BookingController{
		Log:            log,
		Validate:       validate,
		BookingUsecase: booking_usecase,
	}

	router.GET("/bookings", bc.GetAllBookings)
	router.GET("/booking/:id", bc.GetBookingByID)
	router.GET("/booking/order/:id", bc.GetBookingByOrderID)
	router.GET("/booking/payment/callback", bc.GetBookingByID)

	protected.POST("/booking/create", bc.CreateBooking)
	protected.PUT("/booking/update/:id", bc.UpdateBooking)
	protected.DELETE("/booking/:id", bc.DeleteBooking)
}

func (bc *BookingController) CreateBooking(ctx *gin.Context) {
	request := new(model.WriteBookingRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		bc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := bc.Validate.Struct(request); err != nil {
		bc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	err := bc.BookingUsecase.CreateBooking(ctx, request)

	if err != nil {
		bc.Log.WithError(err).Error("failed to create booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Booking created successfully", nil))
}

func (bc *BookingController) GetAllBookings(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := bc.BookingUsecase.ListBookings(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		bc.Log.WithError(err).Error("failed to retrieve bookings")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve bookings", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Bookings retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

// GetBookingByID retrieves a booking by its ID
func (bc *BookingController) GetBookingByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		bc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	data, err := bc.BookingUsecase.GetBookingByID(ctx, uint(id))

	if err != nil {
		bc.Log.WithError(err).WithField("id", id).Error("failed to retrieve booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		bc.Log.WithField("id", id).Warn("booking not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Booking retrieved successfully", nil))
}

// GetBookingByOrderID retrieves a booking by its Order ID
func (bc *BookingController) GetBookingByOrderID(ctx *gin.Context) {
	id := ctx.Param("id")

	data, err := bc.BookingUsecase.GetBookingByOrderID(ctx, id)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			bc.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		bc.Log.WithError(err).WithField("orderID", id).Error("failed to retrieve booking by order ID")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve booking", err.Error()))
		return
	}

	if data == nil {
		bc.Log.WithField("orderID", id).Warn("booking not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Booking not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Booking retrieved successfully", nil))
}

func (bc *BookingController) UpdateBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		bc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing booking ID", nil))
		return
	}

	request := new(model.UpdateBookingRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		bc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := bc.Validate.Struct(request); err != nil {
		bc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := bc.BookingUsecase.UpdateBooking(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			bc.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		bc.Log.WithError(err).WithField("id", id).Error("failed to update booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking updated successfully", nil))
}

func (bc *BookingController) DeleteBooking(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		bc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid booking ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid booking ID", err.Error()))
		return
	}

	if err := bc.BookingUsecase.DeleteBooking(ctx, uint(id)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			bc.Log.WithField("id", id).Warn("booking not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("booking not found", nil))
			return
		}

		bc.Log.WithError(err).WithField("id", id).Error("failed to delete booking")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete booking", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Booking deleted successfully", nil))
}
