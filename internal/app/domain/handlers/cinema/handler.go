package cinema

import (
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/helpers/response"
	"movies/internal/app/core/validation"
	dto "movies/internal/app/domain/core/dto/cinema"
	service "movies/internal/app/domain/services/cinema"
	"net/http"
	"strconv"

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
	payload, ok := validation.BindAndValidate[dto.CreateCinemaDTO](h.binder, ctx)
	if !ok {
		return
	}
	cinema, err := h.service.Create(ctx, dto.CreateCinemaDTO{
		Name:        payload.Name,
		Description: payload.Description,
		Address: payload.Address,
		Latitude: payload.Latitude,
		Longitude: payload.Longitude,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateCinema service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(cinema, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	CinemaID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteCinema handler error")

		return
	}

	ok, err := h.service.Delete(ctx, CinemaID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteCinema service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	CinemaID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetCinema handler error")

		return
	}

	film, err := h.service.Get(ctx, CinemaID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetCinema service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(film, response.OK))
}

