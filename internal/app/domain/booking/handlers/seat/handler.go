package seat

import (
	"gateway/internal/app/core/helpers/errorhandler"
	services "gateway/internal/app/domain/booking/services/seat"
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

// Create CreateSeat godoc
// @Summary		Create Seat
// @Description	Creates a Seat
// @Tags		Seat
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		request  body	seat.CreateSeatDTO	true	"Seat data"
// @Success		200		{object}	response.CommonResponse
// @Router		/seat/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateSeat of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteSeat godoc
// @Summary		Delete Seat by ID
// @Description	Deletes a Seat by ID
// @Tags		Seat
// @Param		id	path	int64	true	"Seat ID"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/seat/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteSeat of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
