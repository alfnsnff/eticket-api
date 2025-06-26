package controller

import (
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
	sc := &ClaimSessionController{
		Log:                 log,
		Validate:            validate,
		ClaimSessionUsecase: claim_session_usecase,
	}

	router.POST("/session/ticket/lock", sc.CreateClaimSession)
	router.POST("/claim/create", sc.TESTCreateClaimSession)
	router.GET("/sessions", sc.GetAllClaimSessions)
	router.GET("/session/:id", sc.GetSessionByID)
	router.POST("/session/ticket/data/entry", sc.UpdateClaimSession)
	router.GET("/session/uuid/:sessionuuid", sc.GetClaimSessionByUUID)
	router.DELETE("/session/:id", sc.DeleteClaimSession)

}

func (csc *ClaimSessionController) CreateClaimSession(ctx *gin.Context) {
	request := new(model.WriteClaimSessionLockTicketsRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		csc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := csc.Validate.Struct(request); err != nil {
		csc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := csc.ClaimSessionUsecase.CreateClaimSession(ctx, request)

	if err != nil {
		csc.Log.WithError(err).Error("failed to create claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}
	csc.Log.WithField("sessionID", datas.SessionID).Info("Claim session created successfully")
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("session_id", datas.SessionID, 60*60, "/", "", true, true)
	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
}

func (csc *ClaimSessionController) TESTCreateClaimSession(ctx *gin.Context) {
	request := new(model.TESTWriteClaimSessionLockTicketsRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		csc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := csc.Validate.Struct(request); err != nil {
		csc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := csc.ClaimSessionUsecase.TESTCreateClaimSession(ctx, request)

	if err != nil {
		csc.Log.WithError(err).Error("failed to create claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}
	csc.Log.WithField("sessionID", datas.SessionID).Info("Claim session created successfully")
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("session_id", datas.SessionID, 60*60, "/", "", true, true)
	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
}

func (csc *ClaimSessionController) GetAllClaimSessions(ctx *gin.Context) {

	params := response.GetParams(ctx)
	datas, total, err := csc.ClaimSessionUsecase.ListClaimSessions(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
		csc.Log.WithError(err).Error("failed to retrieve claim sessions")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim sessions", err.Error()))
		return
	}

	csc.Log.WithField("count", total).Info("Claim sessions retrieved successfully")
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

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Claim sessions retrieved successfully", total, params.Limit, params.Page))
}

func (csc *ClaimSessionController) GetSessionByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		csc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid claim session ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	data, err := csc.ClaimSessionUsecase.GetClaimSessionByID(ctx, uint(id))

	if err != nil {
		csc.Log.WithError(err).WithField("id", id).Error("failed to retrieve claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		csc.Log.WithField("id", id).Warn("claim session not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}

func (csc *ClaimSessionController) UpdateClaimSession(ctx *gin.Context) {

	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		csc.Log.WithError(err).Error("missing session ID in request")
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
		return
	}

	request := new(model.WriteClaimSessionDataEntryRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		csc.Log.WithError(err).Error("failed to bind JSON request body")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := csc.Validate.Struct(request); err != nil {
		csc.Log.WithError(err).Error("failed to validate request body")
		errors := validator.ParseErrors(err)
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Validation error", errors))
		return
	}

	datas, err := csc.ClaimSessionUsecase.UpdateClaimSession(ctx, request, sessionID)

	if err != nil {
		csc.Log.WithError(err).WithField("sessionID", sessionID).Error("failed to update claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update claim session", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(datas, "Claim session updated successfully", nil))
}

func (csc *ClaimSessionController) DeleteClaimSession(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		csc.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid claim session ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	if err := csc.ClaimSessionUsecase.DeleteClaimSession(ctx, uint(id)); err != nil {
		csc.Log.WithError(err).WithField("id", id).Error("failed to delete claim session")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete claim session", err.Error()))
		return
	}

	csc.Log.WithField("id", id).Info("Claim session deleted successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session deleted successfully", nil))
}

func (csc *ClaimSessionController) GetClaimSessionByUUID(ctx *gin.Context) {
	id := ctx.Param("sessionuuid")

	if id == "" {
		csc.Log.WithField("sessionuuid", id).Error("empty session UUID provided")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", "sessionid is empty"))
		return
	}

	data, err := csc.ClaimSessionUsecase.GetBySessionID(ctx, id)

	if err != nil {
		csc.Log.WithError(err).WithField("sessionuuid", id).Error("failed to retrieve claim session by UUID")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		csc.Log.WithField("sessionuuid", id).Warn("claim session not found")
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	csc.Log.WithField("sessionuuid", id).Info("Claim session retrieved successfully")
	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}
