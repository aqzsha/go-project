package screening

import (
	"booking/internal/app/core/helpers/errorhandler"
	"booking/internal/app/core/helpers/response"
	"booking/internal/app/core/validation"
	dto "booking/internal/app/domain/core/dto/screening"
	service "booking/internal/app/domain/services/screening"
	"net/http"
	"strconv"
	"time"

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
	payload, ok := validation.BindAndValidate[dto.CreateScreeningDTO](h.binder, ctx)
	if !ok {
		return
	}

	var startDate time.Time
	if payload.StartAt != "" {
		parseStart, err := time.Parse("2006-01-02", payload.StartAt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ValidationError))
			errorhandler.FailOnError(err, "failed to parse start date")
			return
		}

		startDate = parseStart
	}

	var endDate time.Time
	if payload.EndAt != "" {
		parseEnd, err := time.Parse("2006-01-02", payload.EndAt)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ValidationError))
			errorhandler.FailOnError(err, "failed to parse end date")
			return
		}

		endDate = parseEnd
	}

	screening, err := h.service.Create(ctx, dto.ServiceCreateScreeningDTO{
		CinemaID: payload.CinemaID,
		HallID:   payload.HallID,
		FilmID:   payload.FilmID,
		Price:    payload.Price,
		Format:   payload.Format,
		Language: payload.Language,
		StartAt:  startDate,
		EndAt:    endDate,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateScreening service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(screening, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	ScreeningID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteScreening handler error")

		return
	}

	ok, err := h.service.Delete(ctx, ScreeningID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteScreening service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	ScreeningID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetScreening handler error")

		return
	}

	screening, err := h.service.Get(ctx, ScreeningID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetScreening service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(screening, response.OK))
}

func (h *Handler) ListByFilm(ctx *gin.Context) {
	FilmID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "ListByFilm handler error")

		return
	}

	screening, err := h.service.ListByFilm(ctx, FilmID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListByFilm service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(screening, response.OK))
}

func (h *Handler) ListByCinema(ctx *gin.Context) {
	CinemaID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "ListByCinema handler error")

		return
	}

	screening, err := h.service.ListByCinema(ctx, CinemaID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListByCinema service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(screening, response.OK))
}

func (h *Handler) ListByHall(ctx *gin.Context) {
	HallID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "ListByHall handler error")

		return
	}

	screening, err := h.service.ListByHall(ctx, HallID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListByHall service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(screening, response.OK))
}
