package film

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/movies/services/film"
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

// Create CreateMovie godoc
// @Summary		Create Movie
// @Description	Creates a Movie
// @Tags		Movie
// @Accept		json
// @Produce		json
// @Param		request  body	film.CreateFilmDTO	true	"Film data"
// @Success		200		{object}	response.CommonResponse
// @Router		/film/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateFilm of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteMovie godoc
// @Summary		Delete Movie by ID
// @Description	Deletes a Movie by ID
// @Tags		Movie
// @Param		id	path	int64	true	"Movie ID"
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Router		/film/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteFilm of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}


// Get GetMovie godoc
// @Summary		Get Movie by ID
// @Description	Returns Movie details by its ID
// @Tags		Movie
// @Produce		json
// @Param		id	path	int64	true	"Movie ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/film/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetFilm of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// List   ListMovie godoc
// @Summary		Get list of movie
// @Description	Returns list of Movie
// @Tags		Movie
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Router		/film/list [get]
func (h *Handler) List(ctx *gin.Context) {
	resp := h.service.List(ctx.Request.Context())
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListFilm of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}


// CreateGenre CreateGenre godoc
// @Summary		Create Genre
// @Description	Creates a Genre
// @Tags		Movie
// @Accept		json
// @Produce		json
// @Param		request  body	film.CreateFilmGenreDTO	true	"Film Genre data"
// @Success		200		{object}	response.CommonResponse
// @Router		/film/genre/store [post]
func (h *Handler) CreateGenre(ctx *gin.Context) {
	resp := h.service.CreateGenre(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateGenre of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// DeleteGenre DeleteGenre godoc
// @Summary		Delete Genre by ID
// @Description	Deletes a Genre by ID
// @Tags		Movie
// @Param		id	path	int64	true	"Genre ID"
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Router		/film/genre/delete/{id} [delete]
func (h *Handler) DeleteGenre(ctx *gin.Context) {
	resp := h.service.DeleteGenre(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteGenre of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// GetGenreList   ListGenre godoc
// @Summary		Get list of genre
// @Description	Returns list of genre
// @Tags		Movie
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Router		/film/genre/list [get]
func (h *Handler) GetGenreList(ctx *gin.Context) {
	resp := h.service.GetGenreList(ctx.Request.Context())
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetGenreList of Movie service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
