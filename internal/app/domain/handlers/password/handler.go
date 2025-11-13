package password

import (
	"auth/internal/app/core/helpers/errorhandler"
	"auth/internal/app/core/helpers/response"
	"auth/internal/app/core/validation"
	dto "auth/internal/app/domain/core/dto/requests/password"
	serviceDto "auth/internal/app/domain/core/dto/services/password"
	service "auth/internal/app/domain/services/password"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errStatusMap = map[error]int{
	service.ErrUnauthenticated: http.StatusUnauthorized,
	service.ErrNotFound:        http.StatusNotFound,
	service.ErrInvalidToken:    http.StatusUnauthorized,
	service.ErrUserNotFound:    http.StatusNotFound,
	service.ErrTokenExpired:    http.StatusUnauthorized,
	service.ErrInvalidUser:     http.StatusUnauthorized,
	service.ErrInvalidPassword: http.StatusBadRequest,
}

type Handler struct {
	service service.Service
	binder  *validation.Binder
}

func NewHandler(service service.Service, binder *validation.Binder) *Handler {
	return &Handler{service: service, binder: binder}
}

func (h *Handler) ForgotPassword(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.ForgotPasswordDTO](h.binder, ctx)
	if !ok {
		return
	}
	err := h.service.ForgotPassword(ctx, serviceDto.ForgotPasswordDTO{
		Email: payload.Email,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ForgotPassword service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(nil, response.OK))
}

func (h *Handler) ResetPassword(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.ResetPasswordDTO](h.binder, ctx)
	if !ok {
		return
	}
	err := h.service.ResetPassword(ctx, serviceDto.ResetPasswordDTO{
		Token:       payload.Token,
		NewPassword: payload.NewPassword,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ResetPassword service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(nil, response.OK))
}

func (h *Handler) ChangePassword(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.ChangePasswordDTO](h.binder, ctx)
	if !ok {
		return
	}
	token := ctx.GetHeader("Authorization")

	err := h.service.ChangePassword(ctx, serviceDto.ChangePasswordDTO{
		Password:        payload.Password,
		NewPassword:     payload.NewPassword,
		ConfirmPassword: payload.ConfirmPassword,
		Token:           token,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ChangePassword service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(nil, response.OK))
}
