package api

func (s *Server) initDomainMovieRoutes() {
	s.initFilmRoutes()
	s.initCinemaRoutes()
	s.initHallRoutes()
	s.initGenreRoutes()
	s.initFilmReviewRoutes()
}

func (s *Server) initFilmRoutes() {
	FilmRoutes := mainRouter.Group("film/")
	FilmRoutes.POST("/store", s.deps.filmHandler.Create)
	FilmRoutes.GET("/get/:id", s.deps.filmHandler.Get)
	FilmRoutes.DELETE("/delete/:id", s.deps.filmHandler.Delete)
	FilmRoutes.GET("/list", s.deps.filmHandler.List)
	FilmRoutes.POST("/genre/store", s.deps.filmHandler.CreateGenre)
	FilmRoutes.GET("/genre/list", s.deps.filmHandler.GetGenreList)
	FilmRoutes.GET("/genre/delete/:id", s.deps.filmHandler.DeleteGenre)
}

func (s *Server) initCinemaRoutes() {
	CinemaRoutes := mainRouter.Group("cinema/")
	CinemaRoutes.POST("/store", s.deps.cinemaHandler.Create)
	CinemaRoutes.GET("/get/:id", s.deps.cinemaHandler.Get)
	CinemaRoutes.DELETE("/delete/:id", s.deps.cinemaHandler.Delete)
	CinemaRoutes.GET("/list", s.deps.cinemaHandler.List)
}

func (s *Server) initHallRoutes() {
	HallRoutes := mainRouter.Group("hall/")
	HallRoutes.POST("/store", s.deps.hallHandler.Create)
	HallRoutes.GET("/get/:id", s.deps.hallHandler.Get)
	HallRoutes.DELETE("/delete/:id", s.deps.hallHandler.Delete)
	HallRoutes.GET("/list", s.deps.hallHandler.List)
}

func (s *Server) initGenreRoutes() {
	GenreRoutes := authMiddlewareRouter.Group("genre/")
	GenreRoutes.POST("/store", s.deps.genreHandler.Create)
	GenreRoutes.GET("/get/:id", s.deps.genreHandler.Get)
	GenreRoutes.DELETE("/delete/:id", s.deps.genreHandler.Delete)
}

func (s *Server) initFilmReviewRoutes() {
	FilmReviewRoutes := authMiddlewareRouter.Group("/film/review")
	FilmReviewRoutes.POST("/store", s.deps.filmReviewHandler.Create)
	FilmReviewRoutes.GET("/get/:id", s.deps.filmReviewHandler.Get)
	FilmReviewRoutes.DELETE("/delete/:id", s.deps.filmReviewHandler.Delete)
	FilmReviewRoutes.GET("/list", s.deps.filmReviewHandler.List)
	FilmReviewRoutes.GET("/list/:filmId", s.deps.filmReviewHandler.ListFilm)
}
