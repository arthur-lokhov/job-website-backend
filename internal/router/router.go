package router

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"job-website-backend/internal/handlers"
	"job-website-backend/internal/middleware"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Публичные эндпоинты
	r.Get("/vacancies", handlers.GetVacancies)
	r.Get("/vacancies/{id}", handlers.GetVacancyByID)
	r.Get("/vacancies/{id}/form", handlers.GetVacancyForm)
	r.Post("/applications", handlers.CreateApplication)

	// Приватные эндпоинты (middleware.AuthRequired)
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthRequired)
		r.Get("/vacancies", handlers.AdminGetVacancies)
		r.Post("/vacancies", handlers.AdminCreateVacancy)
		r.Patch("/vacancies/{id}", handlers.AdminUpdateVacancy)
		r.Delete("/vacancies/{id}", handlers.AdminDeleteVacancy)
		r.Get("/applications", handlers.AdminGetApplications)
		r.Patch("/applications/{id}", handlers.AdminUpdateApplication)
		r.Delete("/applications/{id}", handlers.AdminDeleteApplication)
	})

	return r
} 