package hall

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/movies/services/hall"
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

// Create CreateHall godoc
// @Summary		Create Hall
// @Description	Creates a Hall
// @Tags		Hall
// @Accept		json
// @Produce		json
// @Security    BearerAuth
// @Param		request  body	hall.CreateHallDTO	true	"Hall data"
// @Success		200	{object}	response.CommonResponse
// @Router		/hall/store [post]
func (h *Handler) Create(ctx *gin.Context) {
	resp := h.service.Create(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "CreateHall of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Get GetHall godoc
// @Summary		Get Hall by ID
// @Tags		Hall
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Hall ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/hall/get/{id} [get]
func (h *Handler) Get(ctx *gin.Context) {
	resp := h.service.Get(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "GetHall of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Delete DeleteHall godoc
// @Summary		Delete Hall by ID
// @Tags		Hall
// @Produce		json
// @Security    BearerAuth
// @Param		id	path	int64	true	"Hall ID"
// @Success		200	{object}	response.CommonResponse
// @Router		/hall/delete/{id} [delete]
func (h *Handler) Delete(ctx *gin.Context) {
	resp := h.service.Delete(ctx.Request.Context(), ctx.Param("id"))
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "DeleteHall of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// List ListHall godoc
// @Summary		Get list of halls
// @Tags		Hall
// @Produce		json
// @Security    BearerAuth
// @Success		200	{object}	response.CommonResponse
// @Router		/hall/list [get]
func (h *Handler) List(ctx *gin.Context) {
	resp := h.service.List(ctx.Request.Context())
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ListHall of Movie service error")
		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
