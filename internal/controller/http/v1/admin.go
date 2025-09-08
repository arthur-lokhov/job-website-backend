// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"job-website-backend/internal/dto"
	"job-website-backend/internal/services"
	"job-website-backend/pkg/validate"

	"github.com/go-chi/chi/v5"
)

type AdminHandler struct {
	service *services.AdminService
	log     *slog.Logger
}

func NewAdminHandler(service *services.AdminService, log *slog.Logger) *AdminHandler {
	return &AdminHandler{
		service: service,
		log:     log,
	}
}

func (h *AdminHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

// AdminGetVacancies gets a list of all vacancies.
// @Summary Get a list of all vacancies
// @Description Get a list of all vacancies
// @Tags admin
// @Accept  json
// @Produce  json
// @Param isActive query bool false "filter by active vacancies"
// @Success 200 {array} AdminVacancyItem
// @Router /admin/vacancies [get]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminGetVacancies(w http.ResponseWriter, r *http.Request) {
	isActive := false
	if isActiveParam := r.URL.Query().Get("isActive"); isActiveParam == "true" {
		isActive = true
	}

	vacancies, err := h.service.GetVacancies(r.Context(), isActive)
	if err != nil {
		h.log.Error("failed to get vacancies", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	var resp []AdminVacancyItem
	for _, v := range vacancies {
		resp = append(resp, AdminVacancyItem{
			ID:         v.ID.String(),
			Name:       v.Name,
			Location:   &LocationDTO{ID: v.Location.ID.String(), Name: v.Location.Name},
			Department: &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.DepartmentName},
			Level:      &LevelDTO{ID: v.Level.ID.String(), Name: v.Level.Name},
			Important:  v.Important,
			Priority:   v.Priority,
			Info:       v.Info,
			IsActive:   v.IsActive,
			CreatedAt:  v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}

// AdminCreateVacancy creates a new vacancy.
// @Summary Create a new vacancy
// @Description Create a new vacancy
// @Tags admin
// @Accept  json
// @Produce  json
// @Param vacancy body VacancyCreateRequest true "Vacancy data"
// @Success 201 {object} CreatedResponse
// @Router /admin/vacancies [post]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminCreateVacancy(w http.ResponseWriter, r *http.Request) {
	var req VacancyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInvalidRequestData}, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Error("failed to validate request", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	createDTO := dto.VacancyCreateDTO{
		Name:         req.Name,
		LocationID:   req.LocationID,
		DepartmentID: req.DepartmentID,
		LevelID:      req.LevelID,
		Important:    req.Important,
		Priority:     req.Priority,
		Info:         req.Info,
		IsActive:     req.IsActive,
	}
	if req.ApplicationForm != nil {
		createDTO.ApplicationForm = &dto.ApplicationFormDTO{
			FormSchema: req.ApplicationForm.FormSchema,
		}
	}

	vacancyID, createdAt, err := h.service.CreateVacancy(r.Context(), createDTO)
	if err != nil {
		h.log.Error("failed to create vacancy", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, r, CreatedResponse{Success: true, ID: vacancyID, CreatedAt: createdAt.Format(time.RFC3339)}, http.StatusCreated)
}

// AdminUpdateVacancy updates a vacancy.
// @Summary Update a vacancy
// @Description Update a vacancy
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "Vacancy ID"
// @Param vacancy body VacancyUpdateRequest true "Vacancy data"
// @Success 200 {object} SuccessResponse
// @Router /admin/vacancies/{id} [patch]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminUpdateVacancy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req VacancyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInvalidRequestData}, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Error("failed to validate request", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	updateDTO := dto.VacancyUpdateDTO{}
	if req.Name != nil {
		updateDTO.Name = *req.Name
	}
	if req.LocationID != nil {
		updateDTO.LocationID = *req.LocationID
	}
	if req.DepartmentID != nil {
		updateDTO.DepartmentID = *req.DepartmentID
	}
	if req.LevelID != nil {
		updateDTO.LevelID = *req.LevelID
	}
	updateDTO.Important = req.Important
	updateDTO.Priority = req.Priority
	if req.Info != nil {
		updateDTO.Info = *req.Info
	}
	if req.ApplicationForm != nil {
		updateDTO.ApplicationForm = &dto.ApplicationFormDTO{
			FormSchema: req.ApplicationForm.FormSchema,
		}
	}
	updateDTO.IsActive = req.IsActive

	vacancyID, err := h.service.UpdateVacancy(r.Context(), id, updateDTO)
	if err != nil {
		h.log.Error("failed to update vacancy", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, r, SuccessResponse{Success: true, ID: vacancyID}, http.StatusOK)
}

// AdminDeleteVacancy deletes a vacancy.
// @Summary Delete a vacancy
// @Description Delete a vacancy
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "Vacancy ID"
// @Success 204 "No Content"
// @Router /admin/vacancies/{id} [delete]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminDeleteVacancy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteVacancy(r.Context(), id)
	if err != nil {
		h.log.Error("failed to delete vacancy", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AdminGetApplications gets a list of all applications.
// @Summary Get a list of all applications
// @Description Get a list of all applications
// @Tags admin
// @Accept  json
// @Produce  json
// @Param status query string false "filter by status"
// @Param vacancyId query string false "filter by vacancy ID"
// @Param search query string false "search by applicant's full name"
// @Success 200 {array} ApplicationItem
// @Router /admin/applications [get]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminGetApplications(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	vacancyID := r.URL.Query().Get("vacancyId")
	search := r.URL.Query().Get("search") // New search parameter

	applications, err := h.service.AdminGetApplications(r.Context(), status, vacancyID, search) // Pass search parameter
	if err != nil {
		h.log.Error("failed to get applications", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	var resp []ApplicationItem
	for _, app := range applications {
		var answeredAt *string
		if app.AnsweredAt != nil {
			s := app.AnsweredAt.Format("2006-01-02T15:04:05Z07:00")
			answeredAt = &s
		}

		resp = append(resp, ApplicationItem{
			ID:        app.ID.String(),
			VacancyID: app.VacancyID.String(),
			// VacancyName is not populated in service layer yet
			Answers: func() map[string]interface{} {
				var ans map[string]interface{}
				_ = json.Unmarshal(app.Answer, &ans)
				return ans
			}(),
			Status:     app.Status,
			CreatedAt:  app.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  app.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			AnsweredAt: answeredAt,
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}

// AdminUpdateApplication updates an application.
// @Summary Update an application
// @Description Update an application
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Param status body ApplicationStatusUpdate true "Application status"
// @Success 200 {object} SuccessResponse
// @Router /admin/applications/{id} [patch]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminUpdateApplication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ApplicationStatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInvalidRequestData}, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Error("failed to validate request", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	appID, err := h.service.UpdateApplicationStatus(r.Context(), id, req.Status)
	if err != nil {
		h.log.Error("failed to update application status", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, r, SuccessResponse{Success: true, ID: appID}, http.StatusOK)
}

// AdminDeleteApplication deletes an application.
// @Summary Delete an application
// @Description Delete an application
// @Tags admin
// @Accept  json
// @Produce  json
// @Param id path string true "Application ID"
// @Success 204 "No Content"
// @Router /admin/applications/{id} [delete]
// @Security ApiKeyAuth
func (h *AdminHandler) AdminDeleteApplication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.service.DeleteApplication(r.Context(), id)
	if err != nil {
		h.log.Error("failed to delete application", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
