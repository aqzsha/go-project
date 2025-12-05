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
)

type dependency struct {
	authService           authService.Service
	authHandler           *authHandler.AuthHandler
	passwordService       passwordService.Service
	passwordHandler       *passwordHandler.Handler
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

	return &dependency{
		authService:           serviceAuth,
		authHandler:           authHandler.NewHandler(serviceAuth),
		passwordService:       servicePassword,
		passwordHandler:       passwordHandler.NewHandler(servicePassword),
	}, nil
}
