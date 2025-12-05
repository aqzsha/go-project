package hall

import (
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/helpers/response"
	"movies/internal/app/core/validation"
	dto "movies/internal/app/domain/core/dto/hall"
	service "movies/internal/app/domain/services/hall"
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
	payload, ok := validation.BindAndValidate[dto.CreateHallDTO](h.binder, ctx)
	if !ok {
		return
	}
	genre, err := h.service.Create(ctx, dto.CreateHallDTO{
		CinemaID: payload.CinemaID,
		Name:     payload.Name,
		Seats:    payload.Seats,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateHall service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(genre, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	HallID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteHall handler error")

		return
	}

	ok, err := h.service.Delete(ctx, HallID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteHall service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	HallID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetHall handler error")

		return
	}

	hall, err := h.service.Get(ctx, HallID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetHall service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(hall, response.OK))
}


func (h *Handler) List(ctx *gin.Context) {
	review, err := h.service.List(ctx)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListHall service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(review, response.OK))
}