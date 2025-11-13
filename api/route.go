package api

import (
	oauthHelper "auth/internal/app/core/helpers/oauth"
	"auth/internal/app/core/http/middleware"
	"auth/internal/app/core/validation"
	authHandler "auth/internal/app/domain/handlers"
	passwordHandler "auth/internal/app/domain/handlers/password"
	userHandler "auth/internal/app/domain/handlers/user"
	authService "auth/internal/app/domain/services"
	passwordService "auth/internal/app/domain/services/password"
	tokenService "auth/internal/app/domain/services/token"
	userService "auth/internal/app/domain/services/user"

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
	mainRouter.Use(middleware.ApiKey())
	binder := validation.NewBinder(s.validator)
	manager := s.initManager()
	passwdHandler := oauthHelper.NewPasswordHandler(s.pgdb)
	oauthServer := s.initOAuthServer(manager, passwdHandler.Password)

	serviceToken := tokenService.Service{OauthServer: oauthServer}
	serviceAuth := authService.NewService(s.pgdb, serviceToken, s.tokenStore)
	handlerAuth := authHandler.NewHandler(serviceAuth, binder)
	servicePassword := passwordService.NewService(s.pgdb, serviceToken)
	handlerPassword := passwordHandler.NewHandler(servicePassword, binder)
	serviceUser := userService.NewService(s.pgdb)
	handlerUser := userHandler.NewHandler(serviceUser, binder)
	s.initAuthRoutes(handlerAuth)
	s.initPasswordRoutes(handlerPassword)
	s.initUserRoutes(handlerUser)
}

func (s *Server) initUserRoutes(handler *userHandler.Handler) {
	userRoutes := mainRouter.Group("/user")
	userRoutes.POST("/store", handler.Create)
	userRoutes.DELETE("/delete", handler.Delete)
}

func (s *Server) initAuthRoutes(handler *authHandler.Handler) {
	authRoutes := mainRouter.Group("")
	authRoutes.POST("/login", handler.Login)
	authRoutes.POST("/refresh", handler.RefreshToken)
	authRoutes.POST("/check", handler.CheckToken)
	authRoutes.POST("/logout", handler.Logout)
}

func (s *Server) initPasswordRoutes(handler *passwordHandler.Handler) {
	passwordRoutes := mainRouter.Group("/password")
	passwordRoutes.POST("/forgot", handler.ForgotPassword)
	passwordRoutes.POST("/reset", handler.ResetPassword)
	passwordRoutes.POST("/change", handler.ChangePassword)
}
