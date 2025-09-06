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

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "microservice-a/docs"
)

// @title sensor-microservice-a
// @version 1.0
// @description This is the API documentation for Microservice A (Data Generator)
// @host localhost:8080
// @BasePath /
func main() {
	gen := grpcclient.NewGenerator("localhost:50051", 1*time.Second)
	gen.Start("Temperature", "A", "1")

	e := echo.New()
	h := httpHandler.NewHandler(gen)
	e.POST("/frequency", h.UpdateFrequency)

	//	Swagger UI endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server in a goroutine
	go func() {
		log.Println("Microservice A REST running on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
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
