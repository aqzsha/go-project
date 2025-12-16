package api

func (s *Server) initDomainMovieRoutes() {
	movies := mainRouter.Group("")

	movies.Any("/film/*path", s.deps.movieProxy.Handle)
	movies.Any("/cinema/*path", s.deps.movieProxy.Handle)
	movies.Any("/genre/*path", s.deps.movieProxy.Handle)
	movies.Any("/hall/*path", s.deps.movieProxy.Handle)
}
