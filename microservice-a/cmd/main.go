package main

import (
	"context"
	"log"
	"microservice-a/internal/api/grpcclient"
	httpHandler "microservice-a/internal/api/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "microservice-a/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title sensor-microservice-a
// @version 1.0
// @description This is the API documentation for Microservice A (Data Generator)
// @BasePath /
func main() {
	// Load .env file
	envFile := os.Getenv("ENV_FILE") // pass via CLI: ENV_FILE=../configs/temperature.env
	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	sensorType := getEnv("SENSOR_TYPE", "Humidity")
	ID1 := getEnv("ID1", "A")
	ID2 := getEnv("ID2", "1")
	port := getEnv("PORT", "8080")
	grpcTarget := getEnv("GRPC_TARGET", "localhost:50051")

	if sensorType == "" || ID1 == "" || ID2 == "" || port == "" {
		log.Fatal("Please set SENSOR_TYPE, ID1, ID2, and PORT environment variables")
	}

	gen := grpcclient.NewGenerator(grpcTarget, 1*time.Second)
	//gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	gen.Start(sensorType, ID1, ID2)

	e := echo.New()
	h := httpHandler.NewHandler(gen)
	e.POST("/frequency", h.UpdateFrequency)

	//	Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server in a goroutine
	go func() {
		log.Printf("Microservice A (%s) REST running on :%s\n", sensorType, port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	// Stop generator
	gen.Stop()

	// Shutdown Echo server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server stopped")
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
