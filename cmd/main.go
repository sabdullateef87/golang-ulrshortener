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
	"urlshortener/internal/config"
	"urlshortener/internal/http/routes"
)

func main() {
	log.Println("Starting URL Shortener Application...")

	// Load server configuration
	serverConfig := config.LoadServerConfig()

	// Run database migrations first
	config.RunMigrations()

	// Initialize database connection
	config.Connect()

	// Setup HTTP server with routes
	server := setupServer(serverConfig)

	// Start server in a goroutine
	go func() {
		address := fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port)
		log.Printf("Server starting on http://%s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Setup graceful shutdown
	setupGracefulShutdown(server)

	log.Println("Application started successfully!")
	log.Printf("Server listening on http://%s:%s", serverConfig.Host, serverConfig.Port)
	log.Println("Database Host:", config.DbConfig.DatabaseHost)
	log.Println("Available endpoints:")
	log.Println("  GET  /health        - Health check")
	log.Println("  POST /api/shorten   - Create short URL")
	log.Println("  GET  /{shortCode}   - Redirect to long URL")
	log.Println("Press Ctrl+C to shutdown...")

	// Keep the application running
	select {}
}

func setupServer(serverConfig *config.ServerConfig) *http.Server {
	// Use the new routes package
	handler := routes.SetupRoutes()

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", serverConfig.Host, serverConfig.Port),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func setupGracefulShutdown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown HTTP server
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		// Close database connections
		config.Close()

		log.Println("Application stopped")
		os.Exit(0)
	}()
}
