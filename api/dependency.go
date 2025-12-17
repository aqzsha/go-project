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

	filmHandler "gateway/internal/app/domain/movies/handlers/film"
	filmService "gateway/internal/app/domain/movies/services"
)

type dependency struct {
	// auth
	authService     authService.Service
	authHandler     *authHandler.AuthHandler
	passwordService passwordService.Service
	passwordHandler *passwordHandler.Handler

	// movies (proxy)
	filmHandler *filmHandler.FilmHandler
}

func newDeps(baseHttp *myhttp.ClientBase) (*dependency, error) {

	// ---------- AUTH ----------
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

	// ---------- MOVIES ----------
	moviesClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Movies.BaseURL,
		configs.Config.Microservices.Movies.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("microservices movies client: %w", err)
	}

	moviesRH := microservice.NewRequestHandler(moviesClient)
	moviesSvc := filmService.NewService(moviesRH)

	return &dependency{
		authService:     authSvc,
		authHandler:     authHandler.NewHandler(authSvc),
		passwordService: passwordSvc,
		passwordHandler: passwordHandler.NewHandler(passwordSvc),

		filmHandler: filmHandler.NewHandler(moviesSvc),
	}, nil
}
