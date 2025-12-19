package handlers

import (
	"gateway/internal/app/core/helpers/errorhandler"
	"gateway/internal/app/domain/auth/services"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.Service
}

func NewHandler(service services.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login LoginUser godoc
// @Summary		Login
// @Description	Login
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param       request body auth.LoginDTO true "Login data"
// @Success		200		{object}	response.CommonResponse
// @Router		/auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	resp := h.service.Login(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "Login of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Refresh RefreshToken godoc
// @Summary		RefreshToken
// @Description	Get a refresh token
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		request body auth.RefreshTokenDTO true "Refresh token"
// @Success		200	{object}	response.CommonResponse
// @Router		/auth/refresh [post]
func (h *AuthHandler) Refresh(ctx *gin.Context) {
	resp := h.service.Refresh(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "Refresh of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// Logout LogoutUser godoc
// @Summary		Logout
// @Description	Logout
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Success		200	{object}	response.CommonResponse
// @Security    BearerAuth
// @Router		/auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	resp := h.service.Logout(ctx.Request.Context())
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "Logout of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}
