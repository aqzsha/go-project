package user

import (
	"auth/internal/app/core/helpers/errorhandler"
	"auth/internal/app/core/helpers/response"
	"auth/internal/app/core/validation"
	dto "auth/internal/app/domain/core/dto/requests/user"
	serviceDto "auth/internal/app/domain/core/dto/services/user"
	"auth/internal/app/domain/resources"
	service "auth/internal/app/domain/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errStatusMap = map[error]int{
	service.ErrNotFound: http.StatusNotFound,
}

type Handler struct {
	service service.Service
	binder  *validation.Binder
}

func NewHandler(service service.Service, binder *validation.Binder) *Handler {
	return &Handler{service: service, binder: binder}
}

func (h *Handler) Create(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.CreateDTO](h.binder, ctx)
	if !ok {
		return
	}
	user, err := h.service.Create(ctx, serviceDto.CreateDTO{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateUser service error")

		return
	}

	userResource := resources.NewResource(user)
	ctx.JSON(http.StatusOK, response.SuccessResponse(userResource, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.DeleteDTO](h.binder, ctx)
	if !ok {
		return
	}
	ok, err := h.service.Delete(ctx, serviceDto.DeleteDTO{
		Email: payload.Email,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteUser service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}
