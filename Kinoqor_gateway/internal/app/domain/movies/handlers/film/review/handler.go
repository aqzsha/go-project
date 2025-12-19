package review

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/movies/services/film/review"
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

// Create CreateReview godoc
// @Summary		Create Review
// @Description	Creates a Review
// @Tags		Review
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		request  body	review.CreateReviewDTO	true	"Review data"
// @Success		200		{object}	response.CommonResponse
// @Router		/film/review/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateReview of Review service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteReview godoc
// @Summary		Delete Review by ID
// @Description	Deletes a Review by ID
// @Tags		Review
// @Param		id	path	int64	true	"Review ID"
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/film/review/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteReview of Review service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}


// Get GetReview godoc
// @Summary		Get Review by ID
// @Description	Returns Review details by its ID
// @Tags		Review
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Review ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/film/review/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetReview of Review service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// List   ListReview godoc
// @Summary		Get list of Review
// @Description	Returns list of Review
// @Tags		Review
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/film/review/list [get]
func (h *Handler) List(ctx *gin.Context) {
	resp := h.service.List(ctx.Request.Context())
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListReview of Review service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ListFilm   ListReviewFilm godoc
// @Summary		Get list of Review
// @Description	Returns list of Review
// @Tags		Review
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Movie ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/film/review/list/{id} [get]
func (h *Handler) ListFilm(ctx *gin.Context) {
	resp := h.service.ListFilm(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListReview of Review service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

