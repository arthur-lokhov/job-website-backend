// Package services contains the business logic of the application.
package services

import (
	"context"
	"log/slog"
	"time"

	"job-website-backend/internal/dto"
	"job-website-backend/internal/models"
	"job-website-backend/internal/repository"
)

type AdminService struct {
	repo            *repository.AdminRepository
	applicationRepo *repository.ApplicationRepository // Add ApplicationRepository
	log             *slog.Logger
}

func NewAdminService(repo *repository.AdminRepository, applicationRepo *repository.ApplicationRepository, log *slog.Logger) *AdminService {
	return &AdminService{repo: repo, applicationRepo: applicationRepo, log: log}
}

func (s *AdminService) GetVacancies(ctx context.Context, isActive bool) ([]models.Vacancy, error) {
	const op = "services.AdminService.GetVacancies"
	log := s.log.With(slog.String("op", op), slog.Bool("isActive", isActive))

	log.Info("getting vacancies")

	vacancies, err := s.repo.GetVacancies(ctx, isActive)
	if err != nil {
		log.Error("failed to get vacancies from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got vacancies")
	return vacancies, nil
}

func (s *AdminService) AdminGetApplications(ctx context.Context, status, vacancyID, search string) ([]models.Application, error) {
	const op = "services.AdminService.AdminGetApplications"
	log := s.log.With(slog.String("op", op), slog.String("status", status), slog.String("vacancyID", vacancyID), slog.String("search", search))

	log.Info("getting applications")

	applications, err := s.applicationRepo.GetApplicationsByAdmin(ctx, status, vacancyID, search)
	if err != nil {
		log.Error("failed to get applications from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got applications")
	return applications, nil
}

func (s *AdminService) UpdateApplicationStatus(ctx context.Context, id string, status string) (string, error) {
	const op = "services.AdminService.UpdateApplicationStatus"
	log := s.log.With(slog.String("op", op), slog.String("id", id), slog.String("status", status))

	log.Info("updating application status")

	appID, err := s.applicationRepo.UpdateApplicationStatus(ctx, id, status)
	if err != nil {
		log.Error("failed to update application status in repo", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully updated application status")
	return appID, nil
}

func (s *AdminService) CreateVacancy(ctx context.Context, req dto.VacancyCreateDTO) (string, time.Time, error) {
	const op = "services.AdminService.CreateVacancy"
	log := s.log.With(slog.String("op", op))

	log.Info("creating vacancy")

	vacancyID, createdAt, err := s.repo.CreateVacancy(ctx, req)
	if err != nil {
		log.Error("failed to create vacancy in repo", slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	log.Info("successfully created vacancy")
	return vacancyID, createdAt, nil
}

func (s *AdminService) UpdateVacancy(ctx context.Context, id string, req dto.VacancyUpdateDTO) (string, error) {
	const op = "services.AdminService.UpdateVacancy"
	log := s.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("updating vacancy")

	vacancyID, err := s.repo.UpdateVacancy(ctx, id, req)
	if err != nil {
		log.Error("failed to update vacancy in repo", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully updated vacancy")
	return vacancyID, nil
}

func (s *AdminService) DeleteVacancy(ctx context.Context, id string) error {
	const op = "services.AdminService.DeleteVacancy"
	log := s.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("deleting vacancy")

	err := s.repo.DeleteVacancy(ctx, id)
	if err != nil {
		log.Error("failed to delete vacancy in repo", slog.String("error", err.Error()))
		return err
	}

	log.Info("successfully deleted vacancy")
	return nil
}

func (s *AdminService) DeleteApplication(ctx context.Context, id string) error {
	const op = "services.AdminService.DeleteApplication"
	log := s.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("deleting application")

	err := s.applicationRepo.DeleteApplication(ctx, id)
	if err != nil {
		log.Error("failed to delete application in repo", slog.String("error", err.Error()))
		return err
	}

	log.Info("successfully deleted application")
	return nil
}
