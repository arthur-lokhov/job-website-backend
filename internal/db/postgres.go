package db

import (
	"fmt"
	"log"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/config"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// TODO: Убрать SQL код
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	DB.Exec(`
	DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'application_status') THEN
			CREATE TYPE application_status AS ENUM ('NEW', 'IN_REVIEW', 'ACCEPTED', 'REJECTED');
		END IF;
	END$$;
`)

	err = DB.AutoMigrate(
		&models.Department{},
		&models.Level{},
		&models.Location{},
		&models.Vacancy{},
		&models.Application{},
	)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
}
