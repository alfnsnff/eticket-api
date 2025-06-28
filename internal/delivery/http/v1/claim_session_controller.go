package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	"eticket-api/internal/model" // Import the response package
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClaimSessionController struct {
	Validate            validator.Validator
	Log                 logger.Logger
	ClaimSessionUsecase *usecase.ClaimSessionUsecase
}

func NewClaimSessionController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	claim_session_usecase *usecase.ClaimSessionUsecase,

) {
	c := &ClaimSessionController{
		Log:                 log,
		Validate:            validate,
		ClaimSessionUsecase: claim_session_usecase,
	}

	router.POST("/claim/lock", c.LockClaimSession)
	router.POST("/claim/entry", c.EntryClaimSession)
	router.POST("/claim/create", c.CreateClaimSession)
	router.GET("/claims", c.GetAllClaimSessions)
	router.GET("/claim/:id", c.GetSessionByID)
	router.GET("/claim/session/:sessionid", c.GetClaimSessionByUUID)
	router.PUT("/claim/update/:id", c.UpdateClaimSession)
	router.DELETE("/claim/:id", c.DeleteClaimSession)

}

func (c *ClaimSessionController) LockClaimSession(ctx *gin.Context) {
	request := new(model.TESTWriteClaimSessionRequest)

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

	datas, err := c.ClaimSessionUsecase.LockClaimSession(ctx, request)

	if err != nil {
		c.Log.WithError(err).Error("failed to create claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("session_id", datas.SessionID, 60*60, "/", "", true, true)

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
}

func (c *ClaimSessionController) EntryClaimSession(ctx *gin.Context) {

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		c.Log.WithError(err).Error("missing session ID in request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
		return
	}

	request := new(model.TESTWriteClaimSessionDataEntryRequest)
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

	datas, err := c.ClaimSessionUsecase.EntryClaimSession(ctx, request, sessionID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", sessionID).Warn("claim session not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("claim session not found", nil))
			return
		}

		c.Log.WithError(err).WithField("sessionID", sessionID).Error("failed to update claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update claim session", err.Error()))
		return
	}

	// Clear cookies
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("session_id", "", -1, "/", "", true, true)
	ctx.SetCookie("order_id", datas.OrderID, 60*60, "/", "", true, true)
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Claim session updated successfully", nil))
}

func (c *ClaimSessionController) CreateClaimSession(ctx *gin.Context) {
	request := new(model.TESTWriteClaimSessionRequest)

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

	if err := c.ClaimSessionUsecase.CreateClaimSession(ctx, request); err != nil {

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}
		c.Log.WithError(err).Error("failed to create claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Claim session created successfully", nil))
}

func (c *ClaimSessionController) GetAllClaimSessions(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := c.ClaimSessionUsecase.ListClaimSessions(ctx, params.Limit, params.Offset, params.Sort, params.Search)
	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve claim sessions")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim sessions", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		datas,
		"Claim sessions retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))

}

func (c *ClaimSessionController) GetSessionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid claim session ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	data, err := c.ClaimSessionUsecase.GetClaimSessionByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("claim session not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("claim session not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		c.Log.WithField("id", id).Warn("claim session not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}

func (c *ClaimSessionController) UpdateClaimSession(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("failed to parse claim session ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing claim ID", nil))
		return
	}

	request := new(model.UpdateClaimSessionRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		c.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	request.ID = uint(id)
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	if err := c.ClaimSessionUsecase.UpdateClaimSession(ctx, request); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("claim session not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("claim session not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("claim session already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("claim session already exists", nil))
			return
		}

		c.Log.WithError(err).Error("failed to update claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update claim session", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session updated successfully", nil))
}

func (c *ClaimSessionController) DeleteClaimSession(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid claim session ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	if err := c.ClaimSessionUsecase.DeleteClaimSession(ctx, uint(id)); err != nil {

		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("claim session not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("claim session not found", nil))
			return
		}
		c.Log.WithError(err).WithField("id", id).Error("failed to delete claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete claim session", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session deleted successfully", nil))
}

func (c *ClaimSessionController) GetClaimSessionByUUID(ctx *gin.Context) {
	sessionID := ctx.Param("sessionid")

	if sessionID == "" {
		c.Log.WithField("sessionid", sessionID).Error("empty session UUID provided")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", "sessionid is empty"))
		return
	}

	data, err := c.ClaimSessionUsecase.GetBySessionID(ctx, sessionID)

	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", sessionID).Warn("claim session not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("claim session not found", nil))
			return
		}

		c.Log.WithError(err).WithField("sessionid", sessionID).Error("failed to retrieve claim session by UUID")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		c.Log.WithField("sessionid", sessionID).Warn("claim session not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}
