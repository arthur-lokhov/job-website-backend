// Package services contains the business logic of the application.
package services

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"
	"job-website-backend/internal/repository"

	"github.com/google/uuid"
)

type VacancyService struct {
	repo              *repository.VacancyRepository
	departmentService *DepartmentService
	log               *slog.Logger
}

func NewVacancyService(repo *repository.VacancyRepository, departmentService *DepartmentService, log *slog.Logger) *VacancyService {
	return &VacancyService{repo: repo, departmentService: departmentService, log: log}
}

func (s *VacancyService) GetVacancies(ctx context.Context, important bool, sort string, search string) ([]models.Vacancy, error) {
	const op = "services.VacancyService.GetVacancies"
	log := s.log.With(slog.String("op", op), slog.Bool("important", important), slog.String("sort", sort), slog.String("search", search))

	log.Info("getting vacancies")

	vacancies, err := s.repo.GetVacancies(ctx, important, sort, search)
	if err != nil {
		log.Error("failed to get vacancies from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("getting departments from cache")
	// Fetch departments from cache
	departments, err := s.departmentService.GetDepartments(ctx)
	if err != nil {
		log.Error("failed to get departments from cache", slog.String("error", err.Error()))
		return nil, err
	}

	departmentMap := make(map[uuid.UUID]string)
	for _, dept := range departments {
		departmentMap[dept.ID] = dept.Name
	}

	for i := range vacancies {
		if name, ok := departmentMap[vacancies[i].DepartmentID]; ok {
			vacancies[i].DepartmentName = name
		}
	}

	log.Info("successfully got vacancies")
	return vacancies, nil
}

func (s *VacancyService) GetVacancyByID(ctx context.Context, id string) (*models.Vacancy, error) {
	const op = "services.VacancyService.GetVacancyByID"
	log := s.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("getting vacancy by id")

	vacancy, err := s.repo.GetVacancyByID(ctx, id)
	if err != nil {
		log.Error("failed to get vacancy by id from repo", slog.String("error", err.Error()))
		return nil, err
	}

	if vacancy == nil {
		log.Info("vacancy not found")
		return nil, nil // Vacancy not found
	}

	log.Info("getting departments from cache")
	// Fetch departments from cache
	departments, err := s.departmentService.GetDepartments(ctx)
	if err != nil {
		log.Error("failed to get departments from cache", slog.String("error", err.Error()))
		return nil, err
	}

	for _, dept := range departments {
		if dept.ID == vacancy.DepartmentID {
			vacancy.DepartmentName = dept.Name
			break
		}
	}

	log.Info("successfully got vacancy by id")
	return vacancy, nil
}

func (s *VacancyService) GetVacancyForm(ctx context.Context, id string) ([]byte, error) {
	const op = "services.VacancyService.GetVacancyForm"
	log := s.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("getting vacancy form")

	form, err := s.repo.GetVacancyForm(ctx, id)
	if err != nil {
		log.Error("failed to get vacancy form from repo", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got vacancy form")
	return form, nil
}
