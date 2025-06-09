package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/models"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/services"
)

type Handlers struct {
	services *services.Services
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{services: services}
}

// Public handlers
func (h *Handlers) GetVacancies(w http.ResponseWriter, r *http.Request) {
	important := r.URL.Query().Get("important")
	var importantBool *bool
	if important != "" {
		b, err := strconv.ParseBool(important)
		if err == nil {
			importantBool = &b
		}
	}

	vacancies, err := h.services.Vacancy.List(r.Context(), boolPtr(true), importantBool)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vacancies)
}

func (h *Handlers) GetVacancy(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vacancy ID", http.StatusBadRequest)
		return
	}

	vacancy, err := h.services.Vacancy.Get(r.Context(), id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Vacancy not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vacancy)
}

func (h *Handlers) GetVacancyForm(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vacancy ID", http.StatusBadRequest)
		return
	}

	vacancy, err := h.services.Vacancy.Get(r.Context(), id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Vacancy not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"formSchema": vacancy.ApplicationForm,
	})
}

func (h *Handlers) SubmitApplication(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VacancyID uuid.UUID       `json:"vacancyId"`
		Answers   json.RawMessage `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	application := &models.Application{
		VacancyID: req.VacancyID,
		Answer:    req.Answers,
		Status:    "PENDING",
	}

	if err := h.services.Application.Create(r.Context(), application); err != nil {
		if err == services.ErrInvalidInput {
			http.Error(w, "Invalid application data", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      application.ID,
	})
}

// Admin handlers
func (h *Handlers) AdminGetVacancies(w http.ResponseWriter, r *http.Request) {
	isActive := r.URL.Query().Get("isActive")
	var isActiveBool *bool
	if isActive != "" {
		b, err := strconv.ParseBool(isActive)
		if err == nil {
			isActiveBool = &b
		}
	}

	vacancies, err := h.services.Vacancy.List(r.Context(), isActiveBool, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vacancies)
}

func (h *Handlers) AdminCreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancy
	if err := json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.services.Vacancy.Create(r.Context(), &vacancy); err != nil {
		if err == services.ErrInvalidInput {
			http.Error(w, "Invalid vacancy data", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"id":        vacancy.ID,
		"createdAt": vacancy.CreatedAt,
	})
}

func (h *Handlers) AdminUpdateVacancy(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vacancy ID", http.StatusBadRequest)
		return
	}

	var vacancy models.Vacancy
	if err := json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	vacancy.ID = id

	if err := h.services.Vacancy.Update(r.Context(), &vacancy); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Vacancy not found", http.StatusNotFound)
			return
		}
		if err == services.ErrInvalidInput {
			http.Error(w, "Invalid vacancy data", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      vacancy.ID,
	})
}

func (h *Handlers) AdminDeleteVacancy(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid vacancy ID", http.StatusBadRequest)
		return
	}

	if err := h.services.Vacancy.Delete(r.Context(), id); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Vacancy not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) AdminGetApplications(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	vacancyID := r.URL.Query().Get("vacancyId")
	var vacancyIDPtr *uuid.UUID
	if vacancyID != "" {
		id, err := uuid.Parse(vacancyID)
		if err == nil {
			vacancyIDPtr = &id
		}
	}

	applications, err := h.services.Application.List(r.Context(), statusPtr, vacancyIDPtr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(applications)
}

func (h *Handlers) AdminUpdateApplication(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	application, err := h.services.Application.Get(r.Context(), id)
	if err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	application.Status = req.Status
	if err := h.services.Application.Update(r.Context(), application); err != nil {
		if err == services.ErrInvalidInput {
			http.Error(w, "Invalid application data", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      application.ID,
	})
}

func (h *Handlers) AdminDeleteApplication(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	if err := h.services.Application.Delete(r.Context(), id); err != nil {
		if err == services.ErrNotFound {
			http.Error(w, "Application not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
} 