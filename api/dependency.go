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
	cinemaHandler "gateway/internal/app/domain/movies/handlers/cinema"
	cinemaService "gateway/internal/app/domain/movies/services/cinema"
	hallHandler "gateway/internal/app/domain/movies/handlers/hall"
	hallService "gateway/internal/app/domain/movies/services/hall"
	genreHandler "gateway/internal/app/domain/movies/handlers/genre"
	genreService "gateway/internal/app/domain/movies/services/genre"
)

type dependency struct {
	authService           authService.Service
	authHandler           *authHandler.AuthHandler
	passwordService       passwordService.Service
	passwordHandler       *passwordHandler.Handler
	filmService           filmService.Service
	filmHandler           *filmHandler.Handler
	cinemaService         cinemaService.Service
	cinemaHandler         *cinemaHandler.Handler
	hallService 		  hallService.Service
	hallHandler 		  *hallHandler.Handler
	genreService 		  genreService.Service
	genreHandler 		  *genreHandler.Handler
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
	serviceCinema := cinemaService.NewService(movieRequestHandler)
	serviceHall := hallService.NewService(movieRequestHandler)
	return &dependency{
		authService:           serviceAuth,
		authHandler:           authHandler.NewHandler(serviceAuth),
		passwordService:       servicePassword,
		passwordHandler:       passwordHandler.NewHandler(servicePassword),
		filmService:           serviceFilm,
		filmHandler:           filmHandler.NewHandler(serviceFilm),
		cinemaService: serviceCinema,
		cinemaHandler: cinemaHandler.NewHandler(serviceCinema),
		hallHandler: hallHandler.NewHandler(serviceHall),
		hallService: serviceHall,
	}, nil
}