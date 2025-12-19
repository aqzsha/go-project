package genre

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/movies/services/genre"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

// Create CreateGenre godoc
// @Summary		Create Genre
// @Tags		Genre
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		request  body	genre.CreateGenreDTO	true	"Genre data"
// @Success		200	{object}	response.CommonResponse
// @Router		/genre/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "CreateGenre of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Get GetGenre godoc
// @Summary		Get Genre by ID
// @Tags		Genre
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Genre ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/genre/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "GetGenre of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteGenre godoc
// @Summary		Delete Genre by ID
// @Tags		Genre
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Genre ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/genre/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			response.ErrorResponse(response.ServerError),
		)
		errorhandler.FailOnError(resp.Error, "DeleteGenre of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
