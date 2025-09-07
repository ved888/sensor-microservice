package main

import (
	"context"
	"microservice-b/database"
	"microservice-b/internal/api/grpc"
	httpHandler "microservice-b/internal/api/http"
	"microservice-b/internal/repository"
	"microservice-b/internal/usecase"
	myMiddleware "microservice-b/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "microservice-b/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title sensor-microservice-b
// @version 1.0
// @description This is the API documentation for Microservice B (Data Receiver / API Service)
// @host localhost:8081
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	// Structured logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	// Initialize database connection
	db, err := database.DbConnection()
	if err != nil {
		log.WithError(err).Fatal("database connection failed")
	}
	defer db.Close()

	// Repositories
	sensorRepository := repository.NewSensorRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Usecases
	userUseCase := &usecase.UserRepository{
		Repo:      userRepo,
		JWTSecret: "my-secret-key",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Start gRPC server in goroutine
	go grpc.StartGRPCServer(ctx, sensorRepository, ":50051")
	log.Println("Microservice B started. gRPC server listening on :50051")

	// Start Echo REST server
	e := echo.New()

	// Global Middleware
	e.Use(middleware.Recover()) // recover from panics
	e.Use(middleware.Logger())  // log HTTP requests
	e.Use(middleware.CORS())    // allow cross-origin requests

	// Handlers
	sensorHandler := httpHandler.NewSensorHandler(sensorRepository)
	userHandler := httpHandler.NewUserHandler(userUseCase)

	// Public routes
	e.POST("/signup", userHandler.Signup)
	e.POST("/login", userHandler.Login)

	// Protected routes
	apiGroup := e.Group("/api")
	apiGroup.Use(myMiddleware.JWTMiddleware(userUseCase.JWTSecret))

	apiGroup.GET("/sensors", sensorHandler.GetSensors)
	apiGroup.DELETE("/sensors", sensorHandler.DeleteSensors)
	apiGroup.PATCH("/sensors", sensorHandler.EditSensors)

	// Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Run Echo server in goroutine
	go func() {
		log.Println("Microservice B REST server running on :8000")
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Cancel context to stop gRPC server
	cancel()

	shutdownCtx, shutDownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutDownCancel()
	// Shutdown Echo server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.WithError(err).Error("REST server shutdown failed")
	}

	// Stop gRPC server if needed
	log.Println("Microservice B stopped")
}
