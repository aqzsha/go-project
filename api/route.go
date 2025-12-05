package api

import (
	"gateway/configs"
	"gateway/internal/app/core/healthcheck"
	"gateway/internal/app/core/http/middleware/authenticate"

	"github.com/gin-gonic/gin"
)

var mainRouter *gin.RouterGroup
var authMiddlewareRouter *gin.RouterGroup

func (s *Server) initRoutes() error {
	handler = router()

	initSwagger(handler)

	if err := s.initHealthCheck(); err != nil {
		return err
	}

	s.initDomainRoutes()

	return nil
}

func (s *Server) initHealthCheck() error {
	handler.GET("/health/check", healthcheck.HealthHandler(healthcheck.Deps{
		Redis: s.rdb,
		Externals: map[string]string{
			"kinoqor-movies":       configs.Config.Microservices.Movies.BaseURL,
			"kinoqor-auth":         configs.Config.Microservices.Auth.BaseURL,
			"kinoqor-booking":      configs.Config.Microservices.Booking.BaseURL,
			"kinoqor-notification": configs.Config.Microservices.Notification.BaseURL,
		},
	}))
	return nil
}

func (s *Server) initDomainRoutes() {
	mainRouter = handler.Group("/")
	authMiddleware := authenticate.NewMiddleware(s.deps.authService)

	authMiddlewareRouter = mainRouter.Group("")
	authMiddlewareRouter.Use(authMiddleware.Handle())

	s.initDomainAuthRoutes()
}