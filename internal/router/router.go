package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/handlers"
	authMiddleware "git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/middleware"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/services"
)

func NewRouter(h *handlers.Handlers, authService *services.AuthService) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/api/v1/vacancies", h.GetVacancies)
		r.Get("/api/v1/vacancies/{id}", h.GetVacancy)
		r.Get("/api/v1/vacancies/{id}/form", h.GetVacancyForm)
		r.Post("/api/v1/applications", h.SubmitApplication)
	})

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Auth(authService))
		r.Use(authMiddleware.RequirePermission("admin"))

		// Vacancy management
		r.Get("/api/v1/admin/vacancies", h.AdminGetVacancies)
		r.Post("/api/v1/admin/vacancies", h.AdminCreateVacancy)
		r.Put("/api/v1/admin/vacancies/{id}", h.AdminUpdateVacancy)
		r.Delete("/api/v1/admin/vacancies/{id}", h.AdminDeleteVacancy)

		// Application management
		r.Get("/api/v1/admin/applications", h.AdminGetApplications)
		r.Put("/api/v1/admin/applications/{id}", h.AdminUpdateApplication)
		r.Delete("/api/v1/admin/applications/{id}", h.AdminDeleteApplication)
	})

	return r
} 