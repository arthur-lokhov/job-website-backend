// Package v1 contains the full implementation of the public API and admin panel
package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"job-website-backend/internal/services"
)

type LevelsHandler struct {
	service *services.LevelService
	log     *slog.Logger
}

func NewLevelsHandler(service *services.LevelService, log *slog.Logger) *LevelsHandler {
	return &LevelsHandler{service: service, log: log}
}

func (h *LevelsHandler) respondJSON(w http.ResponseWriter, r *http.Request, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()), slog.String("path", r.URL.Path))
	}
}

type LevelResponse struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// GetLevels gets a list of levels.
// @Summary Get a list of levels
// @Description Get a list of levels
// @Tags levels
// @Accept  json
// @Produce  json
// @Success 200 {array} LevelResponse
// @Router /levels [get]
func (h *LevelsHandler) GetLevels(w http.ResponseWriter, r *http.Request) {
	levels, err := h.service.GetLevels(r.Context())
	if err != nil {
		h.log.Error("failed to get levels", slog.String("error", err.Error()))
		h.respondJSON(w, r, map[string]string{"error": ErrInternalServerError}, http.StatusInternalServerError)
		return
	}

	var resp []LevelResponse
	for _, l := range levels {
		resp = append(resp, LevelResponse{
			Label: l.Name,
			Value: l.ID.String(),
		})
	}

	h.respondJSON(w, r, resp, http.StatusOK)
}
