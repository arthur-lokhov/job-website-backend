package main

import (
	"fmt"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/config"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/db"
)

func main() {
	cfg := config.LoadConfig()
	db.Init(cfg)

	fmt.Println("✅ Connected to PostgreSQL via GORM!")
}
