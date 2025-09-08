// Package services contains the business logic of the application.
package services

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"
	"job-website-backend/internal/repository"
)

type LevelService struct {
	repo *repository.LevelRepository
	log  *slog.Logger
}

func NewLevelService(repo *repository.LevelRepository, log *slog.Logger) *LevelService {
	return &LevelService{repo: repo, log: log}
}

func (s *LevelService) GetLevels(ctx context.Context) ([]models.Level, error) {
	const op = "services.LevelService.GetLevels"
	log := s.log.With(slog.String("op", op))

	log.Info("getting levels")

	levels, err := s.repo.GetLevels(ctx)
	if err != nil {
		log.Error("failed to get levels from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got levels")
	return levels, nil
}
