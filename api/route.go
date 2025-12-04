package api

import (
	"movies/internal/app/core/validation"
	cinemaHandler "movies/internal/app/domain/handlers/cinema"
	filmHandler "movies/internal/app/domain/handlers/film"
	genreHandler "movies/internal/app/domain/handlers/genre"
	cinemaService "movies/internal/app/domain/services/cinema"
	filmService "movies/internal/app/domain/services/film"
	genreService "movies/internal/app/domain/services/genre"

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
	
	s.initFilmRoutes(handlerFilm)
	s.initCinemaRoutes(handlerCinema)
	s.initGenreRoutes(handlerGenre)
}

func (s *Server) initFilmRoutes(handler *filmHandler.Handler) {
	FilmRoutes := mainRouter.Group("/film")
	FilmRoutes.POST("/store", handler.Create)
	FilmRoutes.GET("/get/:id", handler.Get)
	FilmRoutes.DELETE("/delete/:id", handler.Delete)
	FilmRoutes.POST("/genre/store", handler.CreateGenre)
	FilmRoutes.GET("/genre/list", handler.GetGenreList)
	FilmRoutes.DELETE("/genre/delete/:id", handler.DeleteGenre)
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
