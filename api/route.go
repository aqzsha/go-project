package api

import (
	"movies/internal/app/core/validation"
	cinemaHandler "movies/internal/app/domain/handlers/cinema"
	filmHandler "movies/internal/app/domain/handlers/film"
	filmReviewHandler "movies/internal/app/domain/handlers/film/review"
	genreHandler "movies/internal/app/domain/handlers/genre"
	hallHandler "movies/internal/app/domain/handlers/hall"
	cinemaService "movies/internal/app/domain/services/cinema"
	filmService "movies/internal/app/domain/services/film"
	filmReviewService "movies/internal/app/domain/services/film/review"
	genreService "movies/internal/app/domain/services/genre"
	hallService "movies/internal/app/domain/services/hall"

	"github.com/gin-gonic/gin"
)

var mainRouter *gin.RouterGroup

func (s *Server) initRoutes() error {
	handler = router()

	if err := s.initHealthCheck(); err != nil {
		return err
	}

	s.initDomainRoutes()

	return nil
}

func (s *Server) initHealthCheck() error {
	return nil
}

func (s *Server) initDomainRoutes() {
	mainRouter = handler.Group("")
	binder := validation.NewBinder(s.validator)
	serviceFilm := filmService.NewService(s.pgdb)
	handlerFilm := filmHandler.NewHandler(serviceFilm, binder)
	serviceCinema := cinemaService.NewService(s.pgdb)	
	handlerCinema := cinemaHandler.NewHandler(serviceCinema, binder)
	serviceGenre := genreService.NewService(s.pgdb)
	handlerGenre := genreHandler.NewHandler(serviceGenre, binder)
	serviceFilmReview := filmReviewService.NewService(s.pgdb)
	handlerFilmReview := filmReviewHandler.NewHandler(serviceFilmReview, binder)
	serviceHall := hallService.NewService(s.pgdb)
	handlerHall := hallHandler.NewHandler(serviceHall, binder)
	
	s.initFilmRoutes(handlerFilm)
	s.initCinemaRoutes(handlerCinema)
	s.initGenreRoutes(handlerGenre)
	s.initFilmReviewRoutes(handlerFilmReview)
	s.initHallRoutes(handlerHall)
}

func (s *Server) initFilmRoutes(handler *filmHandler.Handler) {
	FilmRoutes := mainRouter.Group("/film")
	FilmRoutes.POST("/store", handler.Create)
	FilmRoutes.GET("/get/:id", handler.Get)
	FilmRoutes.DELETE("/delete/:id", handler.Delete)
	FilmRoutes.GET("/list", handler.List)
	FilmRoutes.POST("/genre/store", handler.CreateGenre)
	FilmRoutes.GET("/genre/list", handler.GetGenreList)
	FilmRoutes.DELETE("/genre/delete/:id", handler.DeleteGenre)
}

func (s *Server) initFilmReviewRoutes(handler *filmReviewHandler.Handler){
	FilmReviewRoutes := mainRouter.Group("/film/review")
	FilmReviewRoutes.POST("/store", handler.Create)
	FilmReviewRoutes.GET("/get/:id", handler.Get)
	FilmReviewRoutes.DELETE("/delete/:id", handler.Delete)
	FilmReviewRoutes.GET("/list", handler.List)
}

func (s *Server) initCinemaRoutes(handler *cinemaHandler.Handler) {
	CinemaRoutes := mainRouter.Group("/cinema")
	CinemaRoutes.POST("/store", handler.Create)
	CinemaRoutes.GET("/get/:id", handler.Get)
	CinemaRoutes.DELETE("/delete/:id", handler.Delete)
}

func (s *Server) initGenreRoutes(handler *genreHandler.Handler){
	GenreRoutes := mainRouter.Group("/genre")
	GenreRoutes.POST("/store", handler.Create)
	GenreRoutes.GET("/get/:id", handler.Get)
	GenreRoutes.DELETE("/delete/:id", handler.Delete)
}

func (s *Server) initHallRoutes(handler *hallHandler.Handler){
	HallRoutes := mainRouter.Group("/hall")
	HallRoutes.POST("/store", handler.Create)
	HallRoutes.GET("/get/:id", handler.Get)
	HallRoutes.DELETE("/delete/:id", handler.Delete)
	HallRoutes.GET("/list", handler.List)
}