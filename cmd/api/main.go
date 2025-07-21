package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"job-website-backend/internal/models"
	"job-website-backend/internal/router"
	"job-website-backend/internal/handlers"
)

func main() {
	// Получаем строку подключения из env или config
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/jobdb?sslmode=disable"
	}
	// Подключаемся к БД
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// Авто-миграция моделей (на dev)
	db.AutoMigrate(&models.Vacancy{}, &models.Application{}, &models.Department{}, &models.Level{}, &models.Location{}, &models.AuthService{})

	// Передаем подключение к БД в handlers
	handlers.DB = db

	r := router.New()
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 