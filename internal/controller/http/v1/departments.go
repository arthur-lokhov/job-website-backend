// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"job-website-backend/internal/services"

	"github.com/google/uuid"
)

type DepartmentsHandler struct {
	service *services.DepartmentService
	log     *slog.Logger
}

func NewDepartmentsHandler(service *services.DepartmentService, log *slog.Logger) *DepartmentsHandler {
	return &DepartmentsHandler{service: service, log: log}
}

func (h *DepartmentsHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

type DepartmentResponse struct {
	Label string    `json:"label"`
	ID    uuid.UUID `json:"id"`
	Value string    `json:"value"`
}

// GetDepartments gets a list of departments.
// @Summary Get a list of departments
// @Description Get a list of departments
// @Tags departments
// @Accept  json
// @Produce  json
// @Success 200 {array} DepartmentResponse
// @Router /departments [get]
func (h *DepartmentsHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.service.GetDepartments(r.Context())
	if err != nil {
		h.log.Error("failed to get departments", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	var resp []DepartmentResponse
	for _, d := range departments {
		resp = append(resp, DepartmentResponse{
			Label: d.Name,
			ID:    d.ID,
			Value: d.Slug,
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}
