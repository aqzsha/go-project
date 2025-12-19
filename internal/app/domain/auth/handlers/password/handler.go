package password

import (
	"gateway/internal/app/core/helpers/errorhandler"
	service "gateway/internal/app/domain/auth/services/password"
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

// ForgotPassword godoc
// @Summary		Forgot password
// @Description	Creating a new password
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		request  body	    password.ForgotPasswordDTO	true	"Forgot password"
// @Success		200		{object}	response.CommonResponse
// @Router		/auth/password/forgot [post]
func (h *Handler) ForgotPassword(ctx *gin.Context) {
	resp := h.service.ForgotPassword(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ForgotPassword of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ResetPassword godoc
// @Summary		Reset password
// @Description	Reseting a password
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		request   body	    password.ResetPasswordDTO	true	"Reset password"
// @Success		200		{object}	response.CommonResponse
// @Router		/auth/password/reset [post]
func (h *Handler) ResetPassword(ctx *gin.Context) {
	resp := h.service.ResetPassword(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ResetPassword of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change a password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request    body  password.ChangePasswordDTO  true  "Change password"
// @Success      200   {object}  response.CommonResponse
// @Security     BearerAuth
// @Router       /auth/password/change [post]
func (h *Handler) ChangePassword(ctx *gin.Context) {
	resp := h.service.ChangePassword(ctx.Request.Context(), ctx.Request.Body)
	if resp.Error != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(response.ServerError))
		errorhandler.FailOnError(resp.Error, "ChangePassword of Auth service error")

		return
	}

	ctx.JSON(resp.StatusCode, resp.Data)
}