package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/config"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/db"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/vacancy"
)

func main() {
	cfg := config.LoadConfig()

	db.Init(cfg)

	r := mux.NewRouter()

	// Vacancy endpoints
	r.HandleFunc("/api/vacancies", vacancy.CreateVacancy).Methods("POST")
	r.HandleFunc("/api/vacancies", vacancy.GetAllVacancies).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
