package book

import (
	"booking/internal/app/core/helpers/errorhandler"
	"booking/internal/app/core/helpers/response"
	"booking/internal/app/core/validation"
	service "booking/internal/app/domain/services/book"
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
	ScreeningID, err := strconv.ParseInt(ctx.Param("screeningId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "CreateBooking handler error")

		return
	}

	ok, err := h.service.Create(ctx, ScreeningID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateBooking service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	TicketID, err := strconv.ParseInt(ctx.Param("ticketId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteBooking handler error")

		return
	}

	ok, err := h.service.Delete(ctx, TicketID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteBooking service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	TicketID, err := strconv.ParseInt(ctx.Param("ticketId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetBooking handler error")

		return
	}

	ticket, err := h.service.Get(ctx, TicketID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetBooking service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ticket, response.OK))
}

func (h *Handler) PaidTicket(ctx *gin.Context) {
	seatID, err := strconv.ParseInt(ctx.Param("screeningId"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse("invalid id"))
		return
	}

	pdf, err := h.service.PaidTicket(ctx, seatID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=ticket.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdf)
}

func (h *Handler) ScanTicket(ctx *gin.Context) {
	token := ctx.Param("token")

	err := h.service.ScanTicket(ctx, token)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ScanTicket service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(true, response.OK))
}
