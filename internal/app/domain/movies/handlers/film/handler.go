package film

import (
	film "gateway/internal/app/domain/movies/services"
	"gateway/internal/app/core/helpers/errorhandler"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	service film.Service
}

func NewHandler(service film.Service) *FilmHandler {
	return &FilmHandler{service: service}
}

func (h *FilmHandler) Create(ctx *gin.Context) {
	resp := h.service.Create(
		ctx.Request.Context(),
		ctx.Request.Body,
	)

	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateFilm proxy error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

func (h *FilmHandler) Get(ctx *gin.Context) {
	resp := h.service.Get(
		ctx.Request.Context(),
		ctx.Param("id"),
	)

	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetFilm proxy error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

func (h *FilmHandler) List(ctx *gin.Context) {
	resp := h.service.List(ctx.Request.Context())

	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListFilm proxy error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

func (h *FilmHandler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(
		ctx.Request.Context(),
		ctx.Param("id"),
	)

	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteFilm proxy error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
