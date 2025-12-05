package api

import (
	"context"
	"fmt"
	"gateway/configs"
	"gateway/docs"
	"gateway/internal/app/core/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway/internal/app/core/helpers/errorhandler"
	"gateway/pkg/cache"

	myhttp "gateway/internal/app/core/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	rdb   *redis.Client
	cache cache.Cache
	deps  *dependency
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
	s.rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", configs.Config.Redis.Host, configs.Config.Redis.Port),
	})
	s.cache = cache.NewRedisCache(s.rdb, "gateway")
	baseHttp := myhttp.New(myhttp.ClientConfig{
		Timeout: time.Duration(configs.Config.Microservices.Http.Timeout) * time.Second,
	})
	deps, err := newDeps(baseHttp)
	if err != nil {
		return fmt.Errorf("failed to create dependencies: %w", err)
	}
	s.deps = deps

	return s.initRoutes()
}

func router() *gin.Engine {
	if configs.Config.App.Environment != "local" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Content-Length",
			"X-Requested-With",
			"Accept",
			"Origin",
			"X-CSRF-Token",
			"Cache-Control",
			"Pragma",
			"X-Session-Id",
			"X-api-key",
		},
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
		fmt.Println("kinoqor-gateway - запущен - signal: " + serv.String())
	case err := <-httpServer.Notify():
		fmt.Println(fmt.Errorf("kinoqor-gateway - запущен - httpServer.Notify: %w", err))
	}

	if err := httpServer.Shutdown(); err != nil {
		fmt.Println(fmt.Errorf("kinoqor-gateway - запущен - httpServer.Shutdown: %w", err))
	}

	return nil
}

func initSwagger(r *gin.Engine) {
	docs.SwaggerInfo.Title = "Gateway API"
	docs.SwaggerInfo.Description = "API Gateway"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
