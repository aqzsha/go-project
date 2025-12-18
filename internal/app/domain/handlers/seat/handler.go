package seat

import (
	"booking/internal/app/core/helpers/errorhandler"
	"booking/internal/app/core/helpers/response"
	"booking/internal/app/core/validation"
	dto "booking/internal/app/domain/core/dto/seat"
	service "booking/internal/app/domain/services/seat"
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
	payload, ok := validation.BindAndValidate[dto.CreateSeatDTO](h.binder, ctx)
	if !ok {
		return
	}
	seat, err := h.service.Create(ctx, dto.CreateSeatDTO{
		HallID: payload.HallID,
		Row:    payload.Row,
		Number: payload.Number,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateSeat service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(seat, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	seatID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteSeat handler error")

		return
	}

	ok, err := h.service.Delete(ctx, seatID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteSeat service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}
