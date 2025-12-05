package handler

import (
	"net/http"
	"strconv"

	"screening-service/internal/domain"
	"screening-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ScreeningHandler struct {
	screeningService service.ScreeningService
	seatService      service.SeatService
}

func NewScreeningHandler(screeningService service.ScreeningService, seatService service.SeatService) *ScreeningHandler {
	return &ScreeningHandler{
		screeningService: screeningService,
		seatService:      seatService,
	}
}

// CreateScreening godoc
// @Summary Create a new screening
// @Tags screening
// @Accept json
// @Produce json
// @Param screening body domain.CreateScreeningRequest true "Screening data"
// @Success 201 {object} domain.Screening
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /screenings [post]
func (h *ScreeningHandler) CreateScreening(c *gin.Context) {
	var req domain.CreateScreeningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	screening, err := h.screeningService.CreateScreening(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, screening)
}

// GetScreening godoc
// @Summary Get screening by ID with seat information
// @Tags screening
// @Produce json
// @Param id path int true "Screening ID"
// @Success 200 {object} domain.ScreeningWithSeats
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /screenings/{id} [get]
func (h *ScreeningHandler) GetScreening(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid screening ID"})
		return
	}

	screening, err := h.screeningService.GetScreeningByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, screening)
}

// ListScreenings godoc
// @Summary List screenings with filters
// @Tags screening
// @Produce json
// @Param cinema_id query int false "Cinema ID"
// @Param film_id query int false "Film ID"
// @Param hall_id query int false "Hall ID"
// @Param start_date query string false "Start date (RFC3339)"
// @Param end_date query string false "End date (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} ListScreeningsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /screenings [get]
func (h *ScreeningHandler) ListScreenings(c *gin.Context) {
	var query domain.ListScreeningsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	screenings, total, err := h.screeningService.ListScreenings(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, ListScreeningsResponse{
		Screenings: screenings,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
	})
}

// DeleteScreening godoc
// @Summary Delete a screening
// @Tags screening
// @Param id path int true "Screening ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /screenings/{id} [delete]
func (h *ScreeningHandler) DeleteScreening(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid screening ID"})
		return
	}

	if err := h.screeningService.DeleteScreening(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateSeat godoc
// @Summary Create a seat for a screening
// @Tags seats
// @Accept json
// @Produce json
// @Param seat body domain.CreateSeatRequest true "Seat data"
// @Success 201 {object} domain.Seat
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /seats [post]
func (h *ScreeningHandler) CreateSeat(c *gin.Context) {
	var req domain.CreateSeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	seat, err := h.seatService.CreateSeat(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, seat)
}

// DeleteSeat godoc
// @Summary Delete a seat
// @Tags seats
// @Param id path int true "Seat ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /seats/{id} [delete]
func (h *ScreeningHandler) DeleteSeat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid seat ID"})
		return
	}

	if err := h.seatService.DeleteSeat(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ListScreeningsResponse struct {
	Screenings []domain.Screening `json:"screenings"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
}

func (h *ScreeningHandler) RegisterRoutes(r *gin.RouterGroup) {
	screenings := r.Group("/screenings")
	{
		screenings.POST("", h.CreateScreening)
		screenings.GET("", h.ListScreenings)
		screenings.GET("/:id", h.GetScreening)
		screenings.DELETE("/:id", h.DeleteScreening)
	}

	seats := r.Group("/seats")
	{
		seats.POST("", h.CreateSeat)
		seats.DELETE("/:id", h.DeleteSeat)
	}
}