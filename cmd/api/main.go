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

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/handlers"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/repository"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/router"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/services"
)

// @title Job Website API
// @version 1.0
// @description This is a job website API server.
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Initialize database connection
	dbURL := viper.GetString("database.url")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repositories
	repos := repository.NewRepository(pool)

	// Initialize auth service
	authService := services.NewAuthService(viper.GetString("auth.public_key_url"))

	// Initialize services
	svcs := services.NewServices(repos, authService)

	// Initialize handlers
	h := handlers.NewHandlers(svcs)

	// Initialize router
	r := router.NewRouter(h, authService)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", viper.GetInt("server.port")),
		Handler:      r,
		ReadTimeout:  time.Duration(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("server.writeTimeout")) * time.Second,
		IdleTimeout:  time.Duration(viper.GetInt("server.idleTimeout")) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d", viper.GetInt("server.port"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
} 