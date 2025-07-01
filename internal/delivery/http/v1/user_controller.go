package v1

import (
	"errors"
	errs "eticket-api/internal/common/errors"
	"eticket-api/internal/common/logger"
	"eticket-api/internal/common/validator"
	"eticket-api/internal/delivery/http/response"
	requests "eticket-api/internal/delivery/http/v1/request"
	"eticket-api/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Log         logger.Logger
	Validate    validator.Validator
	UserUsecase *usecase.UserUsecase
}

// NewUserController creates a new UserController instance
func NewUserController(
	router *gin.RouterGroup,
	protected *gin.RouterGroup,
	log logger.Logger,
	validate validator.Validator,
	user_usecase *usecase.UserUsecase,

) {
	c := &UserController{
		Log:         log,
		Validate:    validate,
		UserUsecase: user_usecase,
	}

	router.GET("/users", c.GetAllUsers)
	router.GET("/user/:id", c.GetUserByID)
	router.POST("/user/create", c.CreateUser)

	protected.PUT("/user/update/:id", c.UpdateUser)
	protected.DELETE("/user/:id", c.DeleteUser)
}

func (c *UserController) CreateUser(ctx *gin.Context) {

	request := new(requests.CreateUserRequest)

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

	if err := c.UserUsecase.CreateUser(ctx, requests.UserFromCreate(request)); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}

		c.Log.WithError(err).Error("failed to create user")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to create user", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.NewSuccessResponse(nil, "User created successfully", nil))
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	params := response.GetParams(ctx)
	datas, total, err := c.UserUsecase.ListUsers(ctx, params.Limit, params.Offset, params.Sort, params.Search)
	if err != nil {
		c.Log.WithError(err).Error("failed to retrieve users")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve users", err.Error()))
		return
	}

	responses := make([]*requests.UserResponse, len(datas))
	for i, data := range datas {
		responses[i] = requests.UserToResponse(data)
	}

	ctx.JSON(http.StatusOK, response.NewMetaResponse(
		responses,
		"Users retrieved successfully",
		total,
		params.Limit,
		params.Page,
		params.Sort,
		params.Search,
		params.Path,
	))
}

func (c *UserController) GetUserByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid user ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
		return
	}

	data, err := c.UserUsecase.GetUserByID(ctx, uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("user not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("user not found", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to retrieve user")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to retrieve user", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(requests.UserToResponse(data), "User retrieved successfully", nil))
}

func (c *UserController) UpdateUser(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid or missing user ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid or missing user ID", nil))
		return
	}

	request := new(requests.UpdateUserRequest)
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

	if err := c.UserUsecase.UpdateUser(ctx, requests.UserFromUpdate(request)); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			c.Log.WithField("id", id).Warn("user not found")
			ctx.JSON(http.StatusNotFound, response.NewErrorResponse("user not found", nil))
			return
		}

		if errors.Is(err, errs.ErrConflict) {
			c.Log.WithError(err).Error("user already exists")
			ctx.JSON(http.StatusConflict, response.NewErrorResponse("user already exists", nil))
			return
		}

		c.Log.WithError(err).WithField("id", id).Error("failed to update user")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to update user", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User updated successfully", nil))
}

func (c *UserController) DeleteUser(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Log.WithError(err).WithField("id", ctx.Param("id")).Error("invalid user ID")
		ctx.JSON(http.StatusBadRequest, response.NewErrorResponse("Invalid user ID", err.Error()))
		return
	}

	if err := c.UserUsecase.DeleteUser(ctx, uint(id)); err != nil {
		c.Log.WithError(err).WithField("id", id).Error("failed to delete user")
		ctx.JSON(http.StatusInternalServerError, response.NewErrorResponse("Failed to delete user", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.NewSuccessResponse(nil, "User deleted successfully", nil))
}
