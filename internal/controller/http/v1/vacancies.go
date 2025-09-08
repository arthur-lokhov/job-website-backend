// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"job-website-backend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	ErrMissingVacancyID = "missing vacancy ID in URL"
)

type VacanciesHandler struct {
	service *services.VacancyService
	log     *slog.Logger
}

func NewVacanciesHandler(service *services.VacancyService, log *slog.Logger) *VacanciesHandler {
	return &VacanciesHandler{service: service, log: log}
}

func (h *VacanciesHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

// GetVacancies gets a list of active vacancies.
// @Summary Get a list of active vacancies
// @Description Get a list of active vacancies
// @Tags vacancies
// @Accept  json
// @Produce  json
// @Param important query bool false "filter by important vacancies"
// @Param sort query string false "sort by priority (asc/desc)"
// @Param search query string false "search by vacancy name"
// @Success 200 {array} PublicVacancyItem
// @Router /vacancies [get]
func (h *VacanciesHandler) GetVacancies(w http.ResponseWriter, r *http.Request) {
	important := false
	if imp := r.URL.Query().Get("important"); imp != "" {
		if strings.ToLower(imp) == "true" {
			important = true
		}
	}

	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search") // New search parameter

	vacancies, err := h.service.GetVacancies(r.Context(), important, sort, search) // Pass search parameter
	if err != nil {
		h.log.Error("failed to get vacancies", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": "Internal server error"}, http.StatusInternalServerError)
		return
	}

	var resp []PublicVacancyItem
	for _, v := range vacancies {
		var locationDTO *LocationDTO
		if v.Location.ID != uuid.Nil {
			locationDTO = &LocationDTO{ID: v.Location.ID.String(), Name: v.Location.Name}
		}

		var departmentDTO *DepartmentDTO
		if v.DepartmentID != uuid.Nil {
			departmentDTO = &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.DepartmentName}
		}

		var levelDTO *LevelDTO
		if v.Level.ID != uuid.Nil {
			levelDTO = &LevelDTO{ID: v.Level.ID.String(), Name: v.Level.Name}
		}

		resp = append(resp, PublicVacancyItem{
			ID:         v.ID.String(),
			Name:       v.Name,
			Location:   locationDTO,
			Department: departmentDTO,
			Level:      levelDTO,
			Important:  v.Important,
			Priority:   v.Priority,
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}

// GetVacancyByID gets a vacancy by ID.
// @Summary Get a vacancy by ID
// @Description Get a vacancy by ID
// @Tags vacancies
// @Accept  json
// @Produce  json
// @Param id path string true "Vacancy ID"
// @Success 200 {object} PublicVacancyDetails
// @Router /vacancies/{id} [get]
func (h *VacanciesHandler) GetVacancyByID(w http.ResponseWriter, r *http.Request) {
	// Correct way to get a URL parameter from chi
	id := chi.URLParam(r, "id")

	// Check if the parameter is empty
	if id == "" {
		http.Error(w, ErrMissingVacancyID, http.StatusBadRequest)
		return
	}

	v, err := h.service.GetVacancyByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get vacancy by id", slog.String("id", id), slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": "Internal server error"}, http.StatusInternalServerError)
		return
	}

	var locationDTO *LocationDTO
	if v.Location.ID != uuid.Nil {
		locationDTO = &LocationDTO{ID: v.Location.ID.String(), Name: v.Location.Name}
	}

	var departmentDTO *DepartmentDTO
	if v.DepartmentID != uuid.Nil {
		departmentDTO = &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.DepartmentName}
	}

	var levelDTO *LevelDTO
	if v.Level.ID != uuid.Nil {
		levelDTO = &LevelDTO{ID: v.Level.ID.String(), Name: v.Level.Name}
	}

	resp := PublicVacancyDetails{
		PublicVacancyItem: PublicVacancyItem{
			ID:         v.ID.String(),
			Name:       v.Name,
			Location:   locationDTO,
			Department: departmentDTO,
			Level:      levelDTO,
			Important:  v.Important,
			Priority:   v.Priority,
		},
		Info: v.Info,
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}

// GetVacancyForm gets the application form for a vacancy.
// @Summary Get the application form for a vacancy
// @Description Get the application form for a vacancy
// @Tags vacancies
// @Accept  json
// @Produce  json
// @Param id path string true "Vacancy ID"
// @Success 200 {object} ApplicationFormDTO
// @Router /vacancies/{id}/form [get]
func (h *VacanciesHandler) GetVacancyForm(w http.ResponseWriter, r *http.Request) {
	// Correct way to get a URL parameter from chi
	id := chi.URLParam(r, "id")

	// Check if the parameter is empty
	if id == "" {
		http.Error(w, ErrMissingVacancyID, http.StatusBadRequest)
		return
	}

	form, err := h.service.GetVacancyForm(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get vacancy form", slog.String("id", id), slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": "Internal server error"}, http.StatusInternalServerError)
		return
	}

	var formSchema map[string]interface{}
	if len(form) > 0 {
		_ = json.Unmarshal(form, &formSchema)
	} else {
		formSchema = make(map[string]any)
	}
	resp := ApplicationFormDTO{FormSchema: formSchema}

	h.respondJSON(w, r, resp, http.StatusOK)
}
