package authenticate

import (
	"context"
	"gateway/internal/app/core/contracts/headercontract"
	"gateway/internal/app/core/helpers/errorhandler"
	authservice "gateway/internal/app/domain/auth/services"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	authService authservice.Service
}

func NewMiddleware(authService authservice.Service) *Middleware {
	return &Middleware{
		authService: authService,
	}
}

func (m *Middleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(response.Unauthorized))

			return
		}

		stdCtx := ctx.Request.Context()
		stdCtx = context.WithValue(stdCtx, headercontract.BearerTokenKey{}, token)

		authID, err := m.authService.GetAuthorizedID(stdCtx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(response.Unauthorized))
			errorhandler.FailOnError(err, "failed to get authorized id")

			return
		}
		
		stdCtx = context.WithValue(stdCtx, headercontract.AuthUserKey{}, headercontract.AuthUser{
			ID: authID,
		})
		ctx.Request = ctx.Request.WithContext(stdCtx)

		ctx.Next()
	}
}
