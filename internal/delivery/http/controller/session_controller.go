package controller

// import (
// 	"eticket-api/internal/common/response"
// 	"eticket-api/internal/model" // Import the response package
// 	"eticket-api/internal/usecase/claim_session"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type SessionController struct {
// 	SessionUsecase *claim_session.SessionUsecase
// }

// func NewSessionController(session_usecase *claim_session.SessionUsecase) *SessionController {
// 	return &SessionController{SessionUsecase: session_usecase}
// }

// func (shc *SessionController) CreateSession(ctx *gin.Context) {
// 	request := new(model.WriteClaimSessionRequest)

// 	if err := ctx.ShouldBindJSON(request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	if err := shc.SessionUsecase.CreateSession(ctx, request); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "Claim session created successfully", nil))
// }

// func (shc *SessionController) GetAllSessions(ctx *gin.Context) {
// 	params := response.GetParams(ctx)
// 	datas, total, err := shc.SessionUsecase.GetAllSessions(ctx, params.Limit, params.Offset, params.Sort, params.Search)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim sessions", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.NewMetaResponse(
// 		datas,
// 		"Claim sessions retrieved successfully",
// 		total,
// 		params.Limit,
// 		params.Page,
// 		params.Sort,
// 		params.Search,
// 		params.Path,
// 	))

// 	// ctx.JSON(http.StatusOK, response.NewMetaResponse(datas, "Claim sessions retrieved successfully", total, params.Limit, params.Page))
// }

// func (shc *SessionController) GetSessionByID(ctx *gin.Context) {
// 	id, err := strconv.Atoi(ctx.Param("id"))

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
// 		return
// 	}

// 	data, err := shc.SessionUsecase.GetSessionByID(ctx, uint(id))

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
// 		return
// 	}

// 	if data == nil {
// 		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
// }

// func (shc *SessionController) UpdateSession(ctx *gin.Context) {
// 	request := new(model.UpdateClaimSessionRequest)
// 	id, _ := strconv.Atoi(ctx.Param("id"))

// 	if err := ctx.ShouldBindJSON(request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	if id == 0 {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Claim session ID is required", nil))
// 		return
// 	}

// 	if err := shc.SessionUsecase.UpdateSession(ctx, uint(id), request); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update claim session", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session updated successfully", nil))
// }

// func (shc *SessionController) DeleteSession(ctx *gin.Context) {
// 	id, err := strconv.Atoi(ctx.Param("id"))

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", err.Error()))
// 		return
// 	}

// 	if err := shc.SessionUsecase.DeleteSession(ctx, uint(id)); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete claim session", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "Claim session deleted successfully", nil))
// }

// func (shc *SessionController) GetSessionBySessionID(ctx *gin.Context) {
// 	id := ctx.Param("sessionid")

// 	if id == "" {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid claim session ID", "sessionid is empty"))
// 		return
// 	}

// 	data, err := shc.SessionUsecase.GetBySessionID(ctx, id)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve claim session", err.Error()))
// 		return
// 	}

// 	if data == nil {
// 		ctx.JSON(http.StatusNotFound, response.NewErrorResponse("Claim session not found", nil))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, response.NewSuccessResponse(data, "Claim session retrieved successfully", nil))
// }

// func (csc *SessionController) SessionTicketLock(ctx *gin.Context) {
// 	request := new(model.ClaimedSessionLockTicketsRequest)

// 	if err := ctx.ShouldBindJSON(request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	datas, err := csc.SessionUsecase.SessionLockTickets(ctx, request)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
// 		return
// 	}

// 	ctx.SetSameSite(http.SameSiteNoneMode)
// 	ctx.SetCookie("session_id", datas.SessionID, 60*60, "/", "", true, true)

// 	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
// }

// func (csc *SessionController) SessionTicketDataEntry(ctx *gin.Context) {
// 	sessionID, err := ctx.Cookie("session_id")
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, response.NewErrorResponse("Missing session id", err.Error()))
// 		return
// 	}
// 	request := new(model.ClaimedSessionFillPassengerDataRequest)

// 	if err := ctx.ShouldBindJSON(request); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
// 		return
// 	}

// 	datas, err := csc.SessionUsecase.SessionDataEntry(ctx, request, sessionID)

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create claim session", err.Error()))
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Claim session created successfully", nil))
// }
