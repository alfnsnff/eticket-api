package controller

import (
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/middleware"
	"eticket-api/internal/delivery/response"
	"eticket-api/internal/model" // Import the response package
	"eticket-api/internal/usecase/claim_session"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClaimSessionController struct {
	Validate            validator.Validator
	Log                 logger.Logger
	ClaimSessionUsecase *claim_session.ClaimSessionUsecase
	Authenticate        *middleware.AuthenticateMiddleware
	Authorized          *middleware.AuthorizeMiddleware
}

func NewClaimSessionController(
	router *gin.Engine,
	log logger.Logger,
	validate validator.Validator,
	claim_session_usecase *claim_session.ClaimSessionUsecase,
	authtenticate *middleware.AuthenticateMiddleware,
	authorized *middleware.AuthorizeMiddleware,
) {
	sc := &ClaimSessionController{
		Log:                 log,
		Validate:            validate,
		ClaimSessionUsecase: claim_session_usecase,
		Authenticate:        authtenticate,
		Authorized:          authorized,
	}

	public := router.Group("/api/v1") // No middleware
	public.POST("/session/ticket/lock", sc.CreateClaimSession)
	public.GET("/sessions", sc.GetAllClaimSessions)
	public.GET("/session/:id", sc.GetSessionByID)
	public.POST("/session/ticket/data/entry", sc.UpdateClaimSession)
	public.GET("/session/uuid/:sessionuuid", sc.GetClaimSessionByUUID)
	public.DELETE("/session/:id", sc.DeleteClaimSession)

	protected := router.Group("/api/v1") // No middleware
	protected.Use(sc.Authenticate.Set())
	// protected.Use(ac.Authorized.Set())
}

func (csc *ClaimSessionController) CreateClaimSession(ctx *gin.Context) {
	request := new(model.WriteClaimSessionLockTicketsRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("session_id", datas.SessionID, 60*60, "/", "", true, true)

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
}

func (csc *ClaimSessionController) GetAllClaimSessions(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := csc.ClaimSessionUsecase.GetAllClaimSessions(ctx, params.Limit, params.Offset, params.Sort, params.Search)

	if err != nil {
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

	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Claim sessions retrieved successfully", total, params.Limit, params.Page))
}

func (csc *ClaimSessionController) GetSessionByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	data, err := csc.ClaimSessionUsecase.GetClaimSessionByID(ctx, uint(id))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}

func (csc *ClaimSessionController) UpdateClaimSession(ctx *gin.Context) {
	sessionID, err := ctx.Cookie("session_id")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
		return
	}

	request := new(model.WriteClaimSessionDataEntryRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
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
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
}

func (csc *ClaimSessionController) DeleteClaimSession(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
		return
	}

	if err := csc.ClaimSessionUsecase.DeleteClaimSession(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete claim session", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session deleted successfully", nil))
}

func (csc *ClaimSessionController) GetClaimSessionByUUID(ctx *gin.Context) {
	id := ctx.Param("sessionid")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", "sessionid is empty"))
		return
	}

	data, err := csc.ClaimSessionUsecase.GetBySessionID(ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
		return
	}

	if data == nil {
		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
}
