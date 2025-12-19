package seat

import (
	"gateway/internal/app/core/helpers/errorhandler"
	services "gateway/internal/app/domain/booking/services/book"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Service
}

func NewHandler(service services.Service) *Handler {
	return &Handler{service: service}
}

// Create CreateBook godoc
// @Summary		Create Book
// @Description	Creates a Book
// @Tags		Book
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Screening ID"
// @Success		200		{object}	response.CommonResponse
// @Router		/book/store/{id} [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateBook of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// PaidTicket PaidTickets godoc
// @Summary		PaidTicket
// @Description	PaidTicket
// @Tags		Book
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Screening ID"
// @Success		200		{object}	response.CommonResponse
// @Router		/book/paid-ticket/{id} [post]
func (h *Handler) PaidTicket(ctx *gin.Context) {
	resp := h.service.PaidTicket(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "PaidTicket of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteBook godoc
// @Summary		Delete Book by ID
// @Description	Deletes a Book by ID
// @Tags		Book
// @Param		id	path	int64	true	"Ticket ID"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/book/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteBook of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Get GetBook godoc
// @Summary		Get Book by ID
// @Description	Gets a Book by ID
// @Tags		Book
// @Param		id	path	int64	true	"Ticket ID"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/book/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetBook of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ScanTicket ScanTickets godoc
// @Summary		ScanTicket
// @Description	ScanTicket
// @Tags		Book
// @Param		id	path	string	true	"Token"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/book/tickets/scan/{id} [get]
func (h *Handler) ScanTicket(ctx *gin.Context) {
	resp := h.service.ScanTicket(ctx.Request.Context(), ctx.Param("token"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ScanTicket of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
