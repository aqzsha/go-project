package film

import (
	"gateway/internal/app/domain/movies/services"
	 _ "gateway/internal/app/domain/auth/dto/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service film.Service
}

func NewHandler(service film.Service) *Handler {
	return &Handler{service: service}
}

// CreateFilm godoc
// @Summary     Create film
// @Tags        Movies
// @Accept      json
// @Produce     json
// @Param       request body film.CreateFilmDTO true "Film data"
// @Success     200 {object} response.CommonResponse
// @Router      /api/v1/film/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	ctx.JSON(resp.StatusCode, resp.Data)
}

// GetFilm godoc
// @Summary     Get film
// @Tags        Movies
// @Produce     json
// @Param       id path int true "Film ID"
// @Success     200 {object} response.CommonResponse
// @Router      /api/v1/film/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	resp := h.service.Get(ctx.Request.Context(), id)
	ctx.JSON(resp.StatusCode, resp.Data)
}
