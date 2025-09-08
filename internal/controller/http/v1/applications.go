// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"job-website-backend/internal/services"
	"job-website-backend/pkg/validate"
)

type ApplicationsHandler struct {
	service *services.ApplicationService
	log     *slog.Logger
}

func NewApplicationsHandler(service *services.ApplicationService, log *slog.Logger) *ApplicationsHandler {
	return &ApplicationsHandler{service: service, log: log}
}

func (h *ApplicationsHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

// CreateApplication creates a new application.
// @Summary Create a new application
// @Description Create a new application
// @Tags applications
// @Accept  json
// @Produce  json
// @Param application body ApplicationSubmitRequest true "Application data"
// @Success 201 {object} SuccessResponse
// @Router /applications [post]
func (h *ApplicationsHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var req ApplicationSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": "Invalid request data"}, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(req); err != nil {
		h.log.Error("failed to validate request", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	appID, err := h.service.CreateApplication(r.Context(), req.VacancyID, req.Answers)
	if err != nil {
		h.log.Error("failed to create application", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, r, SuccessResponse{Success: true, ID: appID}, http.StatusCreated)
}
