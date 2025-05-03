package controller

import (
	"eticket-api/internal/model"
	"eticket-api/internal/usecase"
	"eticket-api/pkg/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrebookController struct {
	PreBookUsecase *usecase.PrebookTicketsUsecase
}

// NewPrebookController creates a new PrebookController instance
func NewPrebookController(prebook_usecase *usecase.PrebookTicketsUsecase) *PrebookController {
	return &PrebookController{PreBookUsecase: prebook_usecase}
}

func (pc *PrebookController) LockTicket(ctx *gin.Context) {
	request := new(model.ClaimTicketsRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	datas, err := pc.PreBookUsecase.Execute(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create class", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(datas, "Class created successfully", nil))
}
