package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"booking-service/internal/handler"
	"booking-service/internal/repository"
	"booking-service/internal/service"
	"booking-service/internal/worker"
	"booking-service/pkg/cache"
	"booking-service/pkg/database"
	"booking-service/pkg/pdf"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "movie_booking")
	redisHost := getEnv("REDIS_HOST", "localhost:6379")
	port := getEnv("PORT", "8002")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.MigrateBookingTables(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	
	redisCache := cache.NewRedisCache(redisClient)

	ticketRepo := repository.NewTicketRepository(db)

	//MinIO storage
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := getEnv("MINIO_ACCESS_KEY", "minioadmin")
	minioSecretKey := getEnv("MINIO_SECRET_KEY", "minioadmin")
	minioBucket := getEnv("MINIO_BUCKET", "movie-tickets")
	
	minioStorage, err := storage.NewMinIOStorage(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, false)
	if err != nil {
		log.Printf("Warning: MinIO storage initialization failed: %v", err)
	}

	//external clients (mock for now)
	screeningClient := service.NewScreeningClient(getEnv("SCREENING_SERVICE_URL", "http://localhost:8001"))
	
	smtpHost := getEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASSWORD", "")
	fromEmail := getEnv("FROM_EMAIL", "noreply@cinema.com")
	fromName := getEnv("FROM_NAME", "Cinema Booking")
	
	emailService := email.NewEmailService(smtpHost, smtpPort, smtpUser, smtpPass, fromEmail, fromName)
	
	//user client (mock)
	userClient := service.NewUserClient(getEnv("USER_SERVICE_URL", "http://localhost:8000"))
	
	//notification service with RabbitMQ
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	notificationService, err := service.NewNotificationService(emailService, minioStorage, rabbitURL, userClient)
	if err != nil {
		log.Printf("Warning: Notification service initialization failed: %v", err)
	}

	qrCodeGen := qrcode.NewGenerator(256)
	
	pdfGenerator := pdf.NewTicketGenerator("")

	bookingService := service.NewBookingService(
		ticketRepo,
		screeningClient,
		minioStorage,
		redisCache,
		qrCodeGen,
		pdfGenerator,
		notificationService,
	)

	bookingHandler := handler.NewBookingHandler(bookingService)

	router := gin.Default()
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "booking"})
	})

	api := router.Group("/api/v1")
	bookingHandler.RegisterRoutes(api)

	ticketWorker := worker.NewTicketWorker(bookingService)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go ticketWorker.Start(workerCtx)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("Booking service starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	
	workerCancel()
	ticketWorker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}