package handlers

import (
	"auth/internal/app/core/helpers/errorhandler"
	"auth/internal/app/core/helpers/response"
	"auth/internal/app/core/validation"
	dto "auth/internal/app/domain/core/dto/requests"
	serviceDto "auth/internal/app/domain/core/dto/services"
	"auth/internal/app/domain/resources"
	service "auth/internal/app/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errStatusMap = map[error]int{
	service.ErrNotFound:        http.StatusNotFound,
	service.ErrFailedLogin:     http.StatusUnauthorized,
	service.ErrUnauthenticated: http.StatusUnauthorized,
}

type Handler struct {
	service service.Service
	binder  *validation.Binder
}

func NewHandler(service service.Service, binder *validation.Binder) *Handler {
	return &Handler{service: service, binder: binder}
}

func (h *Handler) Login(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.LoginDTO](h.binder, ctx)
	if !ok {
		return
	}

	token, status, err := h.service.Login(ctx, serviceDto.LoginDTO{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "Login service error")

		return
	}

	ctx.JSON(status, response.SuccessResponse(token, response.OK))
}

func (h *Handler) RefreshToken(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.RefreshTokenDTO](h.binder, ctx)
	if !ok {
		return
	}

	token, err := h.service.RefreshToken(ctx, payload.RefreshToken)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "UserRefreshToken service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(token, response.OK))
}

func (h *Handler) Logout(ctx *gin.Context) {
	logout, err := h.service.Logout(ctx, ctx.GetHeader("Authorization"))
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "Logout service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(logout, response.OK))
}

func (h *Handler) CheckToken(ctx *gin.Context) {
	user, err := h.service.CheckToken(ctx, ctx.GetHeader("Authorization"))
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CheckToken service error")

		return
	}

	userResource := resources.NewResource(user)
	ctx.JSON(http.StatusOK, response.SuccessResponse(userResource, response.OK))
}