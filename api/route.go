package api

import (
	"booking/internal/app/core/validation"
	bookHandler "booking/internal/app/domain/handlers/book"
	screeningHandler "booking/internal/app/domain/handlers/screening"
	seatHandler "booking/internal/app/domain/handlers/seat"
	bookService "booking/internal/app/domain/services/book"
	screeningService "booking/internal/app/domain/services/screening"
	seatService "booking/internal/app/domain/services/seat"

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
	serviceScreening := screeningService.NewService(s.pgdb)
	handlerScreening := screeningHandler.NewHandler(serviceScreening, binder)
	serviceSeat := seatService.NewService(s.pgdb)
	handlerSeat := seatHandler.NewHandler(serviceSeat, binder)
	serviceBook := bookService.NewService(s.pgdb, s.cache, s.txm)
	handlerBook := bookHandler.NewHandler(serviceBook, binder)

	s.initScreeningRoutes(handlerScreening)
	s.initSeatRoutes(handlerSeat)
	s.initBookRoutes(handlerBook)
}

func (s *Server) initScreeningRoutes(handler *screeningHandler.Handler) {
	ScreeningRoutes := mainRouter.Group("/screening")
	ScreeningRoutes.POST("/store", handler.Create)
	ScreeningRoutes.GET("/get/:id", handler.Get)
	ScreeningRoutes.DELETE("/delete/:id", handler.Delete)
	ScreeningRoutes.GET("/list-by-film/:filmId", handler.ListByFilm)
	ScreeningRoutes.GET("/list-by-cinema/:cinemaId", handler.ListByCinema)
	ScreeningRoutes.GET("/list-by-hall/:hallId", handler.ListByHall)
}

func (s *Server) initSeatRoutes(handler *seatHandler.Handler) {
	SeatHandler := mainRouter.Group("/seat")
	SeatHandler.POST("/store", handler.Create)
	SeatHandler.DELETE("/delete/:id", handler.Delete)
}

func (s *Server) initBookRoutes(handler *bookHandler.Handler) {
	BookRoutes := mainRouter.Group("/book")
	BookRoutes.POST("/store/:screeningId", handler.Create)
	BookRoutes.POST("/paid-ticket/:screeningId", handler.PaidTicket)
	BookRoutes.GET("/get/:ticketId", handler.Get)
	BookRoutes.DELETE("/delete/:ticketId", handler.Delete)
	mainRouter.GET("/tickets/scan/:token", handler.ScanTicket)
}
