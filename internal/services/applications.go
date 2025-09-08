// Package services contains the business logic of the application.
package services

import (
	"context"
	"encoding/json"
	"log/slog"

	"job-website-backend/internal/repository"

	"github.com/google/uuid"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
	log  *slog.Logger
}

func NewApplicationService(repo *repository.ApplicationRepository, log *slog.Logger) *ApplicationService {
	return &ApplicationService{repo: repo, log: log}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, vacancyID string, answers map[string]interface{}) (string, error) {
	const op = "services.ApplicationService.CreateApplication"
	log := s.log.With(slog.String("op", op), slog.String("vacancyID", vacancyID))

	log.Info("creating application")

	vacID, err := uuid.Parse(vacancyID)
	if err != nil {
		log.Error("failed to parse vacancyID", slog.String("error", err.Error()))
		return "", err
	}

	ans, err := json.Marshal(answers)
	if err != nil {
		log.Error("failed to marshal answers", slog.String("error", err.Error()))
		return "", err
	}

	appID, err := s.repo.CreateApplication(ctx, vacID, ans)
	if err != nil {
		log.Error("failed to create application in repo", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully created application")
	return appID, nil
}
