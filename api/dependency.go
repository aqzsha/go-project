package api

import (
	"fmt"

	"gateway/configs"
	myhttp "gateway/internal/app/core/http"
	microservice "gateway/internal/app/core/microservices"

	authHandler "gateway/internal/app/domain/auth/handlers"
	passwordHandler "gateway/internal/app/domain/auth/handlers/password"
	authService "gateway/internal/app/domain/auth/services"
	passwordService "gateway/internal/app/domain/auth/services/password"

	moviesProxy "gateway/internal/app/domain/movies/handlers"
)

type dependency struct {
	// auth
	authService     authService.Service
	authHandler     *authHandler.AuthHandler
	passwordService passwordService.Service
	passwordHandler *passwordHandler.Handler

	// movies (ТОЛЬКО proxy)
	movieProxy *moviesProxy.ProxyHandler
}

func newDeps(baseHttp *myhttp.ClientBase) (*dependency, error) {

	// --- AUTH ---
	authClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Auth.BaseURL,
		configs.Config.Microservices.Auth.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("microservices auth client: %w", err)
	}

	authRH := microservice.NewRequestHandler(authClient)

	authSvc := authService.NewService(authRH)
	passwordSvc := passwordService.NewService(authRH)

	// --- MOVIES PROXY ---
	movieProxy := moviesProxy.NewProxyHandler(
		baseHttp.Client,
		configs.Config.Microservices.Movies.BaseURL,
	)

	return &dependency{
		authService:     authSvc,
		authHandler:     authHandler.NewHandler(authSvc),
		passwordService: passwordSvc,
		passwordHandler: passwordHandler.NewHandler(passwordSvc),

		movieProxy: movieProxy,
	}, nil
}
