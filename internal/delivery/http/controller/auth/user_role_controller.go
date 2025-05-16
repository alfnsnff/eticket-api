package controller

import (
	authmodel "eticket-api/internal/model/auth"
	authusecase "eticket-api/internal/usecase/auth"
	"eticket-api/pkg/utils/helper/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRoleController struct {
	UserRoleUsecase *authusecase.UserRoleUsecase
}

// NewUserRoleController creates a new UserRoleController instance
func NewUserRoleController(role_usecase *authusecase.UserRoleUsecase) *UserRoleController {
	return &UserRoleController{UserRoleUsecase: role_usecase}
}

func (urc *UserRoleController) CreateUserRole(ctx *gin.Context) {
	request := new(authmodel.WriteUserRoleRequest)

	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := urc.UserRoleUsecase.CreateUserRole(ctx, request); err != nil {
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to assign user role", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "User role assigned successfully", nil))
}
