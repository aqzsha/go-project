package handler

import (
	"net/http"
	"strconv"

	"review-service/internal/domain"
	"review-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview godoc
// @Summary Create a new review for a film
// @Tags reviews
// @Accept json
// @Produce json
// @Param review body domain.CreateReviewRequest true "Review data"
// @Success 201 {object} domain.Review
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Duplicate review"
// @Failure 500 {object} ErrorResponse
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var req domain.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if exists {
		req.UserID = userID.(int)
	}

	review, err := h.reviewService.CreateReview(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrDuplicateReview {
			statusCode = http.StatusConflict
		} else if err == service.ErrUserNotVerified {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReview godoc
// @Summary Get review by ID
// @Tags reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} domain.Review
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid review ID"})
		return
	}

	review, err := h.reviewService.GetReviewByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// ListFilmReviews godoc
// @Summary List reviews for a film
// @Tags reviews
// @Produce json
// @Param film_id path int true "Film ID"
// @Param min_rating query number false "Minimum rating filter"
// @Param max_rating query number false "Maximum rating filter"
// @Param sort_by query string false "Sort by: created_at, rating" default(created_at)
// @Param sort_order query string false "Sort order: asc, desc" default(desc)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ListReviewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /films/{film_id}/reviews [get]
func (h *ReviewHandler) ListFilmReviews(c *gin.Context) {
	filmID, err := strconv.Atoi(c.Param("film_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid film ID"})
		return
	}

	var query domain.ListReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	reviews, total, err := h.reviewService.GetFilmReviews(c.Request.Context(), filmID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListReviewsResponse{
		Reviews:  reviews,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
}

// ListUserReviews godoc
// @Summary List reviews by a user
// @Tags reviews
// @Produce json
// @Param user_id path int true "User ID"
// @Param film_id query int false "Filter by film"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ListReviewsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{user_id}/reviews [get]
func (h *ReviewHandler) ListUserReviews(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
		return
	}

	var query domain.ListReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	reviews, total, err := h.reviewService.GetUserReviews(c.Request.Context(), userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListReviewsResponse{
		Reviews:  reviews,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	})
}

// UpdateReview godoc
// @Summary Update a review
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param review body domain.UpdateReviewRequest true "Updated review data"
// @Success 200 {object} domain.Review
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid review ID"})
		return
	}

	var req domain.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	review, err := h.reviewService.UpdateReview(c.Request.Context(), id, userID.(int), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == service.ErrReviewNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// DeleteReview godoc
// @Summary Delete a review
// @Tags reviews
// @Param id path int true "Review ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reviews/{id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid review ID"})
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.reviewService.DeleteReview(c.Request.Context(), id, userID.(int)); err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == service.ErrReviewNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFilmRatingStats godoc
// @Summary Get rating statistics for a film
// @Tags reviews
// @Produce json
// @Param film_id path int true "Film ID"
// @Success 200 {object} domain.FilmRatingStats
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /films/{film_id}/ratings [get]
func (h *ReviewHandler) GetFilmRatingStats(c *gin.Context) {
	filmID, err := strconv.Atoi(c.Param("film_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid film ID"})
		return
	}

	stats, err := h.reviewService.GetFilmRatingStats(c.Request.Context(), filmID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListReviewsResponse struct {
	Reviews  []domain.Review `json:"reviews"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func (h *ReviewHandler) RegisterRoutes(r *gin.RouterGroup) {
	reviews := r.Group("/reviews")
	{
		reviews.POST("", h.CreateReview)
		reviews.GET("/:id", h.GetReview)
		reviews.PUT("/:id", h.UpdateReview)
		reviews.DELETE("/:id", h.DeleteReview)
	}

	films := r.Group("/films")
	{
		films.GET("/:film_id/reviews", h.ListFilmReviews)
		films.GET("/:film_id/ratings", h.GetFilmRatingStats)
	}

	users := r.Group("/users")
	{
		users.GET("/:user_id/reviews", h.ListUserReviews)
	}
}