package main

import (
	"context"
	"microservice-b/database"
	"microservice-b/internal/api/grpc"
	httpHandler "microservice-b/internal/api/http"
	"microservice-b/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "microservice-b/docs"
)

// @title sensor-microservice-b
// @version 1.0
// @description This is the API documentation for Microservice B (Data Receiver / API Service)
// @host localhost:8081
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

	sensorRepository := repository.NewSensorRepository(db)

	// Start gRPC server in goroutine
	go grpc.StartGRPCServer(sensorRepository, ":50051")
	log.Println("Microservice B started. gRPC server listening on :50051")

	// Start Echo REST server
	e := echo.New()
	handler := httpHandler.NewSensorHandler(sensorRepository)

	e.GET("/sensors", handler.GetSensors)
	e.DELETE("/sensors", handler.DeleteSensors)
	e.PATCH("/sensors", handler.EditSensors)

	// Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Run Echo server in goroutine
	go func() {
		log.Println("Microservice B REST server running on :8081")
		if err := e.Start(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown Echo server
	if err := e.Shutdown(ctx); err != nil {
		log.WithError(err).Error("REST server shutdown failed")
	}

	// Stop gRPC server if needed
	log.Println("Microservice B stopped")
}
