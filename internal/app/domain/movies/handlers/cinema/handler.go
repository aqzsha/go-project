package cinema

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/movies/services/cinema"
	cinemaDto "gateway/internal/app/domain/movies/dto/cinema"

	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}
var _ = cinemaDto.CreateCinemaDTO{}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// Create CreateCinema godoc
// @Summary     Create Cinema
// @Description Creates a Cinema
// @Tags        Cinema
// @Accept      json
// @Produce     json
// @Param request body CreateCinemaDTO true "Cinema data"
// @Success     200 {object} response.CommonResponse
// @Router      /cinema/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "CreateCinema of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteCinema godoc
// @Summary		Delete Cinema by ID
// @Description	Deletes a Cinema by ID
// @Tags		Cinema
// @Param		id	path	int64	true	"Cinema ID"
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Router		/cinema/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "DeleteCinema of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Get GetCinema godoc
// @Summary		Get Cinema by ID
// @Description	Returns Cinema details by its ID
// @Tags		Cinema
// @Produce		json
// @Param		id	path	int64	true	"Cinema ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/cinema/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "GetCinema of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
