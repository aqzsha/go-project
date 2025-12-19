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
	bookHandler "gateway/internal/app/domain/booking/handlers/book"
	screeningHandler "gateway/internal/app/domain/booking/handlers/screening"
	seatHandler "gateway/internal/app/domain/booking/handlers/seat"
	bookService "gateway/internal/app/domain/booking/services/book"
	screeningService "gateway/internal/app/domain/booking/services/screening"
	seatService "gateway/internal/app/domain/booking/services/seat"
	cinemaHandler "gateway/internal/app/domain/movies/handlers/cinema"
	filmHandler "gateway/internal/app/domain/movies/handlers/film"
	filmReviewHandler "gateway/internal/app/domain/movies/handlers/film/review"
	genreHandler "gateway/internal/app/domain/movies/handlers/genre"
	hallHandler "gateway/internal/app/domain/movies/handlers/hall"
	cinemaService "gateway/internal/app/domain/movies/services/cinema"
	filmService "gateway/internal/app/domain/movies/services/film"
	filmReviewService "gateway/internal/app/domain/movies/services/film/review"
	genreService "gateway/internal/app/domain/movies/services/genre"
	hallService "gateway/internal/app/domain/movies/services/hall"
)

type dependency struct {
	authService       authService.Service
	authHandler       *authHandler.AuthHandler
	passwordService   passwordService.Service
	passwordHandler   *passwordHandler.Handler
	filmService       filmService.Service
	filmHandler       *filmHandler.Handler
	cinemaService     cinemaService.Service
	cinemaHandler     *cinemaHandler.Handler
	hallService       hallService.Service
	hallHandler       *hallHandler.Handler
	genreService      genreService.Service
	genreHandler      *genreHandler.Handler
	filmReviewService filmReviewService.Service
	filmReviewHandler *filmReviewHandler.Handler
	screeningService  screeningService.Service
	screeningHandler  *screeningHandler.Handler
	seatService       seatService.Service
	seatHandler       *seatHandler.Handler
	bookService       bookService.Service
	bookHandler       *bookHandler.Handler
}

func newDeps(baseHttp *myhttp.ClientBase) (*dependency, error) {
	authClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Auth.BaseURL,
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
	)
	if err != nil {
		return nil, fmt.Errorf("microservices movie client: %w", err)
	}

	movieRequestHandler := microservice.NewRequestHandler(movieClient)
	serviceFilm := filmService.NewService(movieRequestHandler)
	serviceCinema := cinemaService.NewService(movieRequestHandler)
	serviceHall := hallService.NewService(movieRequestHandler)
	serviceGenre := genreService.NewService(movieRequestHandler)
	serviceFilmReview := filmReviewService.NewService(movieRequestHandler)

	bookingClient, err := microservice.NewBaseClient(
		baseHttp.Client,
		configs.Config.Microservices.Booking.BaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("microservices booking client: %w", err)
	}

	bookingRequestHandler := microservice.NewRequestHandler(bookingClient)
	serviceScreening := screeningService.NewService(bookingRequestHandler)
	serviceSeat := seatService.NewService(bookingRequestHandler)
	serviceBook := bookService.NewService(bookingRequestHandler)

	return &dependency{
		authService:       serviceAuth,
		authHandler:       authHandler.NewHandler(serviceAuth),
		passwordService:   servicePassword,
		passwordHandler:   passwordHandler.NewHandler(servicePassword),
		filmService:       serviceFilm,
		filmHandler:       filmHandler.NewHandler(serviceFilm),
		cinemaService:     serviceCinema,
		cinemaHandler:     cinemaHandler.NewHandler(serviceCinema),
		hallService:       serviceHall,
		hallHandler:       hallHandler.NewHandler(serviceHall),
		genreService:      serviceGenre,
		genreHandler:      genreHandler.NewHandler(serviceGenre),
		filmReviewService: serviceFilmReview,
		filmReviewHandler: filmReviewHandler.NewHandler(serviceFilmReview),
		screeningService:  serviceScreening,
		screeningHandler:  screeningHandler.NewHandler(serviceScreening),
		seatService:       serviceSeat,
		seatHandler:       seatHandler.NewHandler(serviceSeat),
		bookService:       serviceBook,
		bookHandler:       bookHandler.NewHandler(serviceBook),
	}, nil
}
