package controller

import (
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	RoleUsecase *authusecase.RoleUsecase
}

// NewRoleController creates a new RoleController instance
func NewRoleController(role_usecase *authusecase.RoleUsecase) *RoleController {
	return &RoleController{RoleUsecase: role_usecase}
}

func (rc *RoleController) CreateRole(ctx *gin.Context) {
	request := new(authmodel.WriteRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := rc.RoleUsecase.CreateRole(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "role created successfully", nil))
}
