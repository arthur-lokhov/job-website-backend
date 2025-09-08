// Package services contains the business logic of the application.
package services

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"
	"job-website-backend/internal/repository"
)

type LocationService struct {
	repo *repository.LocationRepository
	log  *slog.Logger
}

func NewLocationService(repo *repository.LocationRepository, log *slog.Logger) *LocationService {
	return &LocationService{repo: repo, log: log}
}

func (s *LocationService) GetLocations(ctx context.Context) ([]models.Location, error) {
	const op = "services.LocationService.GetLocations"
	log := s.log.With(slog.String("op", op))

	log.Info("getting locations")

	locations, err := s.repo.GetLocations(ctx)
	if err != nil {
		log.Error("failed to get locations from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got locations")
	return locations, nil
}
