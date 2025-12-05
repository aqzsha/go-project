package review

import (
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/helpers/response"
	"movies/internal/app/core/validation"
	dto "movies/internal/app/domain/core/dto/film/review"
	service "movies/internal/app/domain/services/film/review"
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
	payload, ok := validation.BindAndValidate[dto.CreateReviewDTO](h.binder, ctx)
	if !ok {
		return
	}
	film, err := h.service.Create(ctx, dto.CreateReviewDTO{
		FilmID:        payload.FilmID,
		UserID:        payload.UserID,
		Body:          payload.Body,
		Rating:        payload.Rating,
	})
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "CreateFilmReview service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(film, response.OK))
}

func (h *Handler) Delete(ctx *gin.Context) {
	ReviewID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "DeleteFilmReview handler error")

		return
	}


	ok, err := h.service.Delete(ctx, ReviewID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "DeleteFilmReview service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(ok, response.Deleted))
}

func (h *Handler) Get(ctx *gin.Context) {
	ReviewID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(err, "GetFilmReview handler error")

		return
	}

	review, err := h.service.Get(ctx, ReviewID)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "GetFilmReview service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(review, response.OK))
}


func (h *Handler) List(ctx *gin.Context) {
	review, err := h.service.List(ctx)
	if err != nil {
		if code, ok := errStatusMap[err]; ok {
			ctx.JSON(code, response.ErrorResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		}

		errorhandler.FailOnError(err, "ListFilmReview service error")

		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse(review, response.OK))
}
