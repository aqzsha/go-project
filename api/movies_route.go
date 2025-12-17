package api

func (s *Server) initDomainMovieRoutes() {
	movies := mainRouter.Group("/film")

	movies.POST("", s.deps.filmHandler.Create)
	movies.GET("/:id", s.deps.filmHandler.Get)
	movies.GET("", s.deps.filmHandler.List)
	movies.DELETE("/:id", s.deps.filmHandler.Delete)
}
