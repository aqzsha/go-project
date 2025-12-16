package api

func (s *Server) initDomainMovieRoutes() {

	mainRouter.Any("/film/*path", s.deps.movieProxy.Handle)

	mainRouter.Any("/cinema/*path", s.deps.movieProxy.Handle)

	mainRouter.Any("/genre/*path", s.deps.movieProxy.Handle)

	mainRouter.Any("/hall/*path", s.deps.movieProxy.Handle)
}
