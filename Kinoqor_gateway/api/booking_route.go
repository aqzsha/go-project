package api

func (s *Server) initDomainBookingRoutes() {
	s.initScreeningRoutes()
}

func (s *Server) initScreeningRoutes() {
	ScreeningRoutes := authMiddlewareRouter.Group("screening/")
	ScreeningRoutes.POST("/store", s.deps.screeningHandler.Create)
	ScreeningRoutes.GET("/get/:id", s.deps.screeningHandler.Get)
	ScreeningRoutes.DELETE("/delete/:id", s.deps.screeningHandler.Delete)
	ScreeningRoutes.GET("/list-by-film/:id", s.deps.screeningHandler.ListByFilm)
	ScreeningRoutes.GET("/list-by-cinema/:id", s.deps.screeningHandler.ListByCinema)
	ScreeningRoutes.GET("/list-by-hall/:id", s.deps.screeningHandler.ListByHall)
}

func (s *Server) initSeatRoutes() {
	SeatHandler := authMiddlewareRouter.Group("/seat")
	SeatHandler.POST("/store", s.deps.seatHandler.Create)
	SeatHandler.DELETE("/delete/:id", s.deps.seatHandler.Delete)
}

func (s *Server) initBookRoutes() {
	BookRoutes := authMiddlewareRouter.Group("/book")
	BookRoutes.POST("/store/:id", s.deps.bookHandler.Create)
	BookRoutes.POST("/paid-ticket/:id", s.deps.bookHandler.PaidTicket)
	BookRoutes.GET("/get/:id", s.deps.bookHandler.Get)
	BookRoutes.DELETE("/delete/:id", s.deps.bookHandler.Delete)
	mainRouter.GET("/tickets/scan/:token", s.deps.bookHandler.ScanTicket)
}
