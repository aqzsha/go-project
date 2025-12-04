package api

import (
	"context"
	"fmt"
	"movies/configs"
	"movies/internal/app/core/helpers/errorhandler"
	"movies/internal/app/core/http"
	"movies/internal/app/core/rules"
	"movies/pkg/cache"
	"movies/pkg/postgres"
	"movies/pkg/validation"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	pgdb       *gorm.DB
	rdb        *redis.Client
	cache      cache.Cache
	validator  *validation.Validator
}

var handler *gin.Engine

func NewServer(ctx context.Context) error {
	s := &Server{}

	if err := s.initDeps(ctx); err != nil {
		errorhandler.FailOnError(err, "проблема с инициализацией зависимостей")

		return fmt.Errorf("server initDeps: %w", err)
	}

	return nil
}

func (s *Server) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		s.initConfig,
		s.initLayers,
		s.initServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) initConfig(_ context.Context) error {
	configs.InitConfig()

	return nil
}

func (s *Server) initLayers(_ context.Context) error {
	s.pgdb = postgres.NewClient()

	validator, err := validation.New(validation.UniqueRule{
		Checker: rules.GormChecker{DB: s.pgdb},
		CtxProvider: func() context.Context {
			return context.TODO()
		},
	})
	if err != nil {
		return fmt.Errorf("failed to init validator: %w", err)
	}

	s.validator = validator

	s.rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", configs.Config.Redis.Host, configs.Config.Redis.Port),
	})
	s.cache = cache.NewRedisCache(s.rdb, "users")

	return s.initRoutes()
}

func router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool { return true },
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"Accept",
			"Origin",
			"X-CSRF-Token",
			"Cache-Control",
			"Pragma",
			"X-Session-Id",
			"X-api-key",
		},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Recovery())

	return r
}


func (s *Server) initServer(_ context.Context) error {
	httpCfg := configs.Config.App.Url
	httpServer := http.NewHttpServer(handler, http.Port(httpCfg))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case serv := <-interrupt:
		fmt.Println("kinoqor-movies-api - запущен - signal: " + serv.String())
	case err := <-httpServer.Notify():
		fmt.Println(fmt.Errorf("kinoqor-movies-api - запущен - httpServer.Notify: %w", err))
	}
	
	if err := httpServer.Shutdown(); err != nil {
		fmt.Println(fmt.Errorf("kinoqor-movies-api - запущен - httpServer.Shutdown: %w", err))
	}

	return nil
}
