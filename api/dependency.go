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
	filmService "gateway/internal/app/domain/movies/services/film"
)

type dependency struct {
	authService           authService.Service
	authHandler           *authHandler.AuthHandler
	passwordService       passwordService.Service
	passwordHandler       *passwordHandler.Handler
	filmService           filmService.Service
	filmHandler           *filmHandler.Handler
}

func newDeps(baseHttp *myhttp.ClientBase) (*dependency, error) {
	authClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Auth.BaseURL,
		configs.Config.Microservices.Auth.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("microservices auth client: %w", err)
	}

	authRequestHandler := microservice.NewRequestHandler(authClient)
	serviceAuth := authService.NewService(authRequestHandler)
	servicePassword := passwordService.NewService(authRequestHandler)

	movieClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Movies.BaseURL,
		configs.Config.Microservices.Movies.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("microservices movie client: %w", err)
	}

	movieRequestHandler := microservice.NewRequestHandler(movieClient)
	serviceFilm := filmService.NewService(movieRequestHandler)

	return &dependency{
		authService:           serviceAuth,
		authHandler:           authHandler.NewHandler(serviceAuth),
		passwordService:       servicePassword,
		passwordHandler:       passwordHandler.NewHandler(servicePassword),
		filmService:           serviceFilm,
		filmHandler:           filmHandler.NewHandler(serviceFilm),
	}, nil
}