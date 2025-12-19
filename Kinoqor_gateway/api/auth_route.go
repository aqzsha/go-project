package api

func (s *Server) initDomainAuthRoutes() {
	s.initAuthRoutes()
	s.initPasswordRoutes()
}

func (s *Server) initAuthRoutes() {
	authRoutes := mainRouter.Group("auth/")
	authRoutes.POST("/login", s.deps.authHandler.Login)
	authRoutes.POST("/refresh", s.deps.authHandler.Refresh)
	authMiddlewareRouter.POST("auth/logout", s.deps.authHandler.Logout)
}

func (s *Server) initPasswordRoutes() {
	passwordRoutes := mainRouter.Group("auth/password/")
	passwordRoutes.POST("/forgot", s.deps.passwordHandler.ForgotPassword)
	passwordRoutes.POST("/reset", s.deps.passwordHandler.ResetPassword)
	authMiddlewareRouter.POST("auth/password/change", s.deps.passwordHandler.ChangePassword)
}
