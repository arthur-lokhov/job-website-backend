package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/config"
)

func main() {
	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig(".", logger)
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Подключение к базе данных
	dbpool, err := pgxpool.New(context.Background(), cfg.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbpool.Close()

	r := mux.NewRouter()

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
