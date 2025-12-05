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

		authEmail, err := m.authService.GetAuthorizedEmail(stdCtx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(response.Unauthorized))
			errorhandler.FailOnError(err, "failed to get authorized email")

			return
		}

		authUser, err := m.authService.GetByEmail(stdCtx, authEmail)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(response.Unauthorized))
			errorhandler.FailOnError(err, "failed to get authorized user data")

			return
		}

		stdCtx = context.WithValue(stdCtx, headercontract.AuthUserKey{}, headercontract.AuthUser{
			ID:         authUser.ID,
		})
		ctx.Request = ctx.Request.WithContext(stdCtx)

		ctx.Next()
	}
}
