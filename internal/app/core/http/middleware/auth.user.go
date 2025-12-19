package middleware

import (
	"context"
	headercontract "movies/internal/app/core/contracts/microservices/header-contract"
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/helpers/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ContextWithAuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUser, err := parseAuthUser(ctx)
		if err != nil {
			errMsg := "failed to parse auth user middleware"
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(errMsg))
			errorhandler.FailOnError(err, errMsg)

			return
		}

		auCtx := context.WithValue(ctx.Request.Context(), headercontract.AuthUserKey{}, authUser)
		ctx.Request = ctx.Request.WithContext(auCtx)

		ctx.Next()
	}
}

func parseAuthUser(ctx *gin.Context) (headercontract.AuthUser, error) {
	var authUser headercontract.AuthUser

	id, err := strconv.ParseInt(ctx.GetHeader(headercontract.UserIdKey), 10, 64)
	if err != nil {
		return headercontract.AuthUser{}, err
	}

	authUser.ID = id

	return authUser, nil
}
