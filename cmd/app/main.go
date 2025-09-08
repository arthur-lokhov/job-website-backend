// @title Job Website API
// @version 1.0
// @description This is a sample server for a job website.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "job-website-backend/internal/controller/http/v1"
	"job-website-backend/internal/repository"
	"job-website-backend/internal/router"
	"job-website-backend/internal/services"
	"job-website-backend/pkg/cache"
	"job-website-backend/pkg/config"
	"job-website-backend/pkg/logger"
	"job-website-backend/pkg/postgres"

	_ "github.com/jackc/pgx/v5/stdlib" // Added this line
)

func main() {
	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Logger
	log := logger.New("info")

	log.Info("config and logger initialized")

	// Valkey Client
	valkeyClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Error("failed to connect to Valkey", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := valkeyClient.Close(); err != nil {
			log.Error("failed to close Valkey client", "error", err)
		}
	}()

	// Database
	dbPool, err := postgres.New(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close() // Use dbPool.Close()

	// Migrations
	// Open a separate *sql.DB connection for Goose migrations
	dbMigrate, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		log.Error("failed to open migration database connection", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := dbMigrate.Close(); err != nil {
			log.Error("failed to close migration database connection", "error", err)
		}
	}()

	if err := postgres.RunMigrations(dbMigrate, "migrations"); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Repositories
	vacancyRepo := repository.NewVacancyRepository(dbPool, log)
	applicationRepo := repository.NewApplicationRepository(dbPool, log)
	adminRepo := repository.NewAdminRepository(dbPool, log)
	departmentRepo := repository.NewDepartmentRepository(dbPool, log)
	levelRepo := repository.NewLevelRepository(dbPool, log)
	locationRepo := repository.NewLocationRepository(dbPool, log)

	// Services
	applicationService := services.NewApplicationService(applicationRepo, log)
	adminService := services.NewAdminService(adminRepo, applicationRepo, log)
	departmentService := services.NewDepartmentService(departmentRepo, cfg.AuthURL, valkeyClient, log)
	authService := services.NewAuthService(valkeyClient, log)
	vacancyService := services.NewVacancyService(vacancyRepo, departmentService, log)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			if _, err := departmentService.FetchAndCacheDepartments(); err != nil {
				log.Error("failed to fetch and cache departments", "error", err)
			}
		}
	}()
	levelService := services.NewLevelService(levelRepo, log)
	locationService := services.NewLocationService(locationRepo, log)

	// Handlers
	vacanciesHandler := v1.NewVacanciesHandler(vacancyService, log)
	applicationsHandler := v1.NewApplicationsHandler(applicationService, log)
	adminHandler := v1.NewAdminHandler(adminService, log)
	departmentsHandler := v1.NewDepartmentsHandler(departmentService, log)
	levelsHandler := v1.NewLevelsHandler(levelService, log)
	locationsHandler := v1.NewLocationsHandler(locationService, log)

	// HTTP Server
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Port,
		Handler: router.New(vacanciesHandler, applicationsHandler, adminHandler, departmentsHandler, levelsHandler, locationsHandler, authService, log),
	}

	go func() {
		log.Info("starting server", "port", cfg.HTTPServer.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server exited properly")
}
