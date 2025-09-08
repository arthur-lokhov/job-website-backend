// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"job-website-backend/internal/services"
)

type LocationsHandler struct {
	service *services.LocationService
	log     *slog.Logger
}

func NewLocationsHandler(service *services.LocationService, log *slog.Logger) *LocationsHandler {
	return &LocationsHandler{service: service, log: log}
}

func (h *LocationsHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

type LocationResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// GetLocations gets a list of locations.
// @Summary Get a list of locations
// @Description Get a list of locations
// @Tags locations
// @Accept  json
// @Produce  json
// @Success 200 {array} LocationResponse
// @Router /locations [get]
func (h *LocationsHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.service.GetLocations(r.Context())
	if err != nil {
		h.log.Error("failed to get locations", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	var resp []LocationResponse
	for _, l := range locations {
		resp = append(resp, LocationResponse{
			Label: l.Name,
			Value: l.ID.String(),
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}
