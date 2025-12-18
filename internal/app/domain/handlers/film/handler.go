package film

import (
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/helpers/response"
	"movies/internal/app/core/validation"
	dto "movies/internal/app/domain/core/dto/film"
	service "movies/internal/app/domain/services/film"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var errStatusMap = map[error]int{
	service.ErrNotFound: http.StatusNotFound,
	service.ErrGenreNotFound: http.StatusNotFound,
}

type Handler struct {
	service service.Service
	binder  *validation.Binder
}

func NewHandler(service service.Service, binder *validation.Binder) *Handler {
	return &Handler{service: service, binder: binder}
}


func (h *Handler) Create(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.CreateFilmDTO](h.binder, ctx)
	if !ok {
		return
	}

	var startDate time.Time
	if payload.StartDate != "" {
		parseStart, err := time.Parse("2006-01-02", payload.StartDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ValidationError))
			errorhandler.FailOnError(err, "failed to parse start date")
			return
		}

		startDate = parseStart
	}

	var endDate time.Time
	if payload.StartDate != "" {
		parseEnd, err := time.Parse("2006-01-02", payload.EndDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ValidationError))
			errorhandler.FailOnError(err, "failed to parse end date")
			return
		}

		endDate = parseEnd
	}

	var premier time.Time	
	if payload.StartDate != "" {
		parsePremier, err := time.Parse("2006-01-02", payload.Premier)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse(response.ValidationError))
			errorhandler.FailOnError(err, "failed to parse end date")
			return
		}

		premier = parsePremier
	}



	film, err := h.service.Create(ctx, dto.CreateFilmServiceDTO{
		Name:        payload.Name,
		Description: payload.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		Duration:    payload.Duration,
		Premier:     premier,
		Production:  payload.Production,
		Director:    payload.Director,
		Rate:        payload.Rate,
		AgeLimit:    payload.AgeLimit,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateFilm service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(film, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	FilmID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteFilm handler error")

		return
	}

	ok, err := h.service.Delete(ctx, FilmID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteFilm service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	FilmID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetFilm handler error")

		return
	}

	film, err := h.service.Get(ctx, FilmID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetFilm service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(film, response.OK))
}

func (h *Handler) CreateGenre(ctx *gin.Context) {
	payload, ok := validation.BindAndValidate[dto.CreateFilmGenreDTO](h.binder, ctx)
	if !ok {
		return
	}
	genre, err := h.service.CreateGenre(ctx, dto.CreateFilmGenreDTO{
		FilmID: payload.FilmID,
		GenreID: payload.GenreID,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateFilmGenre service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(genre, response.OK))
}

func (h *Handler) DeleteGenre(ctx *gin.Context) {
	FilmGenreID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteFilmGenre handler error")

		return
	}

	ok, err := h.service.DeleteGenre(ctx, FilmGenreID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteFilmGenre service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}



func (h *Handler) GetGenreList(ctx *gin.Context) {
	film, err := h.service.GetGenreList(ctx)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetGenreList service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(film, response.OK))
}


func (h *Handler) List(ctx *gin.Context) {
	review, err := h.service.List(ctx)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListFilm service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(review, response.OK))
}
