package api

func (s *Server) initDomainMovieRoutes() {
	s.initFilmRoutes()
	s.initCinemaRoutes()
	s.initHallRoutes()
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
}

func (s *Server) initHallRoutes() {
	HallRoutes := mainRouter.Group("hall/")
	HallRoutes.POST("/store", s.deps.hallHandler.Create)
	HallRoutes.GET("/get/:id", s.deps.hallHandler.Get)
	HallRoutes.DELETE("/delete/:id", s.deps.hallHandler.Delete)
	HallRoutes.GET("/list", s.deps.hallHandler.List)
}