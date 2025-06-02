package handlers

import (
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/config"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewRouter(
	cfg *config.Config,
	db *pgxpool.Pool,
	logger *zap.Logger,
	// Хендде
) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/vacancy")
}
