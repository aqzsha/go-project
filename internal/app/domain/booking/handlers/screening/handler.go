package screening

import (
	"gateway/internal/app/core/helpers/errorhandler"
	services "gateway/internal/app/domain/booking/services/screening"
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

// Create CreateScreening godoc
// @Summary		Create Screening
// @Description	Creates a Screening
// @Tags		Screening
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		request  body	screening.CreateScreeningDTO	true	"Screening data"
// @Success		200		{object}	response.CommonResponse
// @Router		/screening/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateScreening of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteScreening godoc
// @Summary		Delete Screening by ID
// @Description	Deletes a Screening by ID
// @Tags		Screening
// @Param		id	path	int64	true	"Screening ID"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/screening/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteScreening of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Get GetScreening godoc
// @Summary		Get Screening by ID
// @Description	Returns Screening details by its ID
// @Tags		Screening
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Screening ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/screening/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetFilm of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ListByFilm   ListScreening godoc
// @Summary		Get list of Screening
// @Description	Returns list of Screening
// @Tags		Screening
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Film ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/screening/list-by-film/{id} [get]
func (h *Handler) ListByFilm(ctx *gin.Context) {
	resp := h.service.ListByFilm(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListScreening of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ListByCinema   ListScreening godoc
// @Summary		Get list of Screening
// @Description	Returns list of Screening
// @Tags		Screening
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Cinema ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/screening/list-by-cinema/{id} [get]
func (h *Handler) ListByCinema(ctx *gin.Context) {
	resp := h.service.ListByCinema(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListScreening of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ListByHall   ListScreening godoc
// @Summary		Get list of Screening
// @Description	Returns list of Screening
// @Tags		Screening
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Hall ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/screening/list-by-hall/{id} [get]
func (h *Handler) ListByHall(ctx *gin.Context) {
	resp := h.service.ListByHall(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListScreening of Booking service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
