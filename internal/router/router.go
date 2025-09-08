// Package router sets up the HTTP routes for the application.
package router

import (
	"log/slog"
	"net/http"

	"job-website-backend/internal/middleware"
	"job-website-backend/internal/services"

	_ "job-website-backend/docs"
	v1 "job-website-backend/internal/controller/http/v1"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	PermissionVacanciesRead      = "jobservice.vacancies:read"
	PermissionVacanciesCreate    = "jobservice.vacancies:create"
	PermissionVacanciesUpdate    = "jobservice.vacancies:update"
	PermissionVacanciesDelete    = "jobservice.vacancies:delete"
	PermissionApplicationsRead   = "jobservice.applications:read"
	PermissionApplicationsUpdate = "jobservice.applications:update"
	PermissionApplicationsDelete = "jobservice.applications:delete"
)

func New(vacanciesHandler *v1.VacanciesHandler, applicationsHandler *v1.ApplicationsHandler, adminHandler *v1.AdminHandler, departmentsHandler *v1.DepartmentsHandler, levelsHandler *v1.LevelsHandler, locationsHandler *v1.LocationsHandler, authService *services.AuthService, log *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Публичные эндпоинты
		r.Get("/vacancies", vacanciesHandler.GetVacancies)
		r.Get("/vacancies/{id}", vacanciesHandler.GetVacancyByID)
		r.Get("/vacancies/{id}/form", vacanciesHandler.GetVacancyForm)
		r.Post("/applications", applicationsHandler.CreateApplication)

		// Departments, Levels, Locations
		r.Get("/departments", departmentsHandler.GetDepartments)
		r.Get("/levels", levelsHandler.GetLevels)
		r.Get("/locations", locationsHandler.GetLocations)
		// Приватные эндпоинты (middleware.AuthRequired)
		r.Route("/admin", func(r chi.Router) {
			r.Get("/vacancies", middleware.AuthRequired(authService, log, PermissionVacanciesRead)(http.HandlerFunc(adminHandler.AdminGetVacancies)))
			r.Post("/vacancies", middleware.AuthRequired(authService, log, PermissionVacanciesCreate)(http.HandlerFunc(adminHandler.AdminCreateVacancy)))
			r.Patch("/vacancies/{id}", middleware.AuthRequired(authService, log, PermissionVacanciesUpdate)(http.HandlerFunc(adminHandler.AdminUpdateVacancy)))
			r.Delete("/vacancies/{id}", middleware.AuthRequired(authService, log, PermissionVacanciesDelete)(http.HandlerFunc(adminHandler.AdminDeleteVacancy)))
			r.Get("/applications", middleware.AuthRequired(authService, log, PermissionApplicationsRead)(http.HandlerFunc(adminHandler.AdminGetApplications)))
			r.Patch("/applications/{id}", middleware.AuthRequired(authService, log, PermissionApplicationsUpdate)(http.HandlerFunc(adminHandler.AdminUpdateApplication)))
			r.Delete("/applications/{id}", middleware.AuthRequired(authService, log, PermissionApplicationsDelete)(http.HandlerFunc(adminHandler.AdminDeleteApplication)))
		})
	})

	return r
}
