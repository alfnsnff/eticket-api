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

type TicketController struct {
	Validate      validator.Validator
	Log           logger.Logger
	TicketUsecase *usecase.TicketUsecase
}

func NewTicketController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	ticket_usecase *usecase.TicketUsecase,
) {
	c := &TicketController{
		Log:           log,
		Validate:      validate,
		TicketUsecase: ticket_usecase,
	}

	router.GET("/tickets", c.GetAllTickets)
	router.GET("/tickets/schedule/:id", c.GetAllTicketsByScheduleID)
	router.GET("/ticket/:id", c.GetTicketByID)

	protected.PATCH("/ticket/check-in/:id", c.CheckIn)
	protected.POST("/ticket/create", c.CreateTicket)
	protected.PUT("/ticket//update:id", c.UpdateTicket)
	protected.DELETE("/ticket/:id", c.DeleteTicket)
}

func (c *TicketController) CreateTicket(ctx *gin.Context) {

	request := new(model.WriteTicketRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	switch request.Type {
	case "passenger":
		if request.IDType == nil || *request.IDType == "" ||
			request.IDNumber == nil || *request.IDNumber == "" {
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", "id_type and id_number are required for passenger ticket"))
			return
		}
	case "vehicle":
		if request.LicensePlate == nil || *request.LicensePlate == "" {
			ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", "license_plate is required for vehicle ticket"))
			return
		}
	}

	if err := c.TicketUsecase.CreateTicket(ctx, request); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}

		c.Log.WithError(err).Error("failed to create ticket")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create ticket", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Ticket created successfully", nil))
}

func (c *TicketController) GetAllTickets(ctx *gin.Context) {
	params := response.GetParams(ctx)

	datas, total, err := c.TicketUsecase.ListTickets(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve tickets")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Tickets retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *TicketController) GetAllTicketsByScheduleID(ctx *gin.Context) {
	params := response.GetParams(ctx)
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse schedule ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid schedule ID", err.Error()))
		return
	}

	datas, total, err := c.TicketUsecase.ListTicketsByScheduleID(ctx, id, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		c.Log.WithError(err).WithField("schedule_id", id).Error("failed to retrieve tickets by schedule ID")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve tickets", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Tickets retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *TicketController) GetTicketByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse ticket ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error()))
		return
	}

	data, err := c.TicketUsecase.GetTicketByID(ctx, uint(id))

	if err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve ticket")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve ticket", err.Error()))
		return
	}

	if data == nil {
		c.Log.WithField("id", id).Warn("ticket not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Ticket not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Ticket retrieved successfully", nil))
}

func (c *TicketController) UpdateTicket(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse ticket ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing ticket ID", nil))
		return
	}

	request := new(model.UpdateTicketRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.TicketUsecase.UpdateTicket(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("ticket not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("ticket not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("ticket already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to update ticket")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update ticket", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket updated successfully", nil))
}

func (c *TicketController) DeleteTicket(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse ticket ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error()))
		return
	}

	if err := c.TicketUsecase.DeleteTicket(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to delete ticket")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete ticket", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket deleted successfully", nil))
}

func (c *TicketController) CheckIn(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse ticket ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid ticket ID", err.Error()))
		return
	}

	if err := c.TicketUsecase.CheckIn(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to check in ticket")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to check in ticket", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Ticket checked in successfully", nil))
}
