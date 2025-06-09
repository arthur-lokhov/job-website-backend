package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/models"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/repository"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInternalError = errors.New("internal error")
)

type Services struct {
	Vacancy     *VacancyService
	Application *ApplicationService
	Department  *DepartmentService
	Level       *LevelService
	Location    *LocationService
	Auth        *AuthService
}

func NewServices(repo repository.Repository, authService *AuthService) *Services {
	return &Services{
		Vacancy:     NewVacancyService(repo),
		Application: NewApplicationService(repo),
		Department:  NewDepartmentService(repo),
		Level:       NewLevelService(repo),
		Location:    NewLocationService(repo),
		Auth:        authService,
	}
}

type VacancyService struct {
	repo repository.Repository
}

func NewVacancyService(repo repository.Repository) *VacancyService {
	return &VacancyService{repo: repo}
}

func (s *VacancyService) Create(ctx context.Context, vacancy *models.Vacancy) error {
	if err := validateVacancy(vacancy); err != nil {
		return err
	}
	return s.repo.CreateVacancy(ctx, vacancy)
}

func (s *VacancyService) Get(ctx context.Context, id uuid.UUID) (*models.Vacancy, error) {
	vacancy, err := s.repo.GetVacancy(ctx, id)
	if err != nil {
		return nil, err
	}
	if vacancy == nil {
		return nil, ErrNotFound
	}
	return vacancy, nil
}

func (s *VacancyService) Update(ctx context.Context, vacancy *models.Vacancy) error {
	if err := validateVacancy(vacancy); err != nil {
		return err
	}
	return s.repo.UpdateVacancy(ctx, vacancy)
}

func (s *VacancyService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteVacancy(ctx, id)
}

func (s *VacancyService) List(ctx context.Context, isActive *bool, important *bool) ([]*models.Vacancy, error) {
	return s.repo.ListVacancies(ctx, isActive, important)
}

func (s *VacancyService) Search(ctx context.Context, query string) ([]*models.Vacancy, error) {
	return s.repo.SearchVacancies(ctx, query)
}

type ApplicationService struct {
	repo repository.Repository
}

func NewApplicationService(repo repository.Repository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (s *ApplicationService) Create(ctx context.Context, application *models.Application) error {
	if err := validateApplication(application); err != nil {
		return err
	}
	application.Status = "PENDING"
	return s.repo.CreateApplication(ctx, application)
}

func (s *ApplicationService) Get(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	application, err := s.repo.GetApplication(ctx, id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, ErrNotFound
	}
	return application, nil
}

func (s *ApplicationService) Update(ctx context.Context, application *models.Application) error {
	if err := validateApplication(application); err != nil {
		return err
	}
	if application.Status != "" {
		now := time.Now()
		application.AnsweredAt = &now
	}
	return s.repo.UpdateApplication(ctx, application)
}

func (s *ApplicationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteApplication(ctx, id)
}

func (s *ApplicationService) List(ctx context.Context, status *string, vacancyID *uuid.UUID) ([]*models.Application, error) {
	return s.repo.ListApplications(ctx, status, vacancyID)
}

type DepartmentService struct {
	repo repository.Repository
}

func NewDepartmentService(repo repository.Repository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) Create(ctx context.Context, department *models.Department) error {
	if err := validateDepartment(department); err != nil {
		return err
	}
	return s.repo.CreateDepartment(ctx, department)
}

func (s *DepartmentService) Get(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	department, err := s.repo.GetDepartment(ctx, id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, ErrNotFound
	}
	return department, nil
}

func (s *DepartmentService) Update(ctx context.Context, department *models.Department) error {
	if err := validateDepartment(department); err != nil {
		return err
	}
	return s.repo.UpdateDepartment(ctx, department)
}

func (s *DepartmentService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDepartment(ctx, id)
}

func (s *DepartmentService) List(ctx context.Context) ([]*models.Department, error) {
	return s.repo.ListDepartments(ctx)
}

type LevelService struct {
	repo repository.Repository
}

func NewLevelService(repo repository.Repository) *LevelService {
	return &LevelService{repo: repo}
}

func (s *LevelService) Create(ctx context.Context, level *models.Level) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	return s.repo.CreateLevel(ctx, level)
}

func (s *LevelService) Get(ctx context.Context, id uuid.UUID) (*models.Level, error) {
	level, err := s.repo.GetLevel(ctx, id)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, ErrNotFound
	}
	return level, nil
}

func (s *LevelService) Update(ctx context.Context, level *models.Level) error {
	if err := validateLevel(level); err != nil {
		return err
	}
	return s.repo.UpdateLevel(ctx, level)
}

func (s *LevelService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLevel(ctx, id)
}

func (s *LevelService) List(ctx context.Context) ([]*models.Level, error) {
	return s.repo.ListLevels(ctx)
}

type LocationService struct {
	repo repository.Repository
}

func NewLocationService(repo repository.Repository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) Create(ctx context.Context, location *models.Location) error {
	if err := validateLocation(location); err != nil {
		return err
	}
	return s.repo.CreateLocation(ctx, location)
}

func (s *LocationService) Get(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	location, err := s.repo.GetLocation(ctx, id)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, ErrNotFound
	}
	return location, nil
}

func (s *LocationService) Update(ctx context.Context, location *models.Location) error {
	if err := validateLocation(location); err != nil {
		return err
	}
	return s.repo.UpdateLocation(ctx, location)
}

func (s *LocationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLocation(ctx, id)
}

func (s *LocationService) List(ctx context.Context) ([]*models.Location, error) {
	return s.repo.ListLocations(ctx)
}

// Validation functions
func validateVacancy(v *models.Vacancy) error {
	if v.Name == "" {
		return ErrInvalidInput
	}
	if v.Info == "" {
		return ErrInvalidInput
	}
	if v.ApplicationForm == nil {
		return ErrInvalidInput
	}
	return nil
}

func validateApplication(a *models.Application) error {
	if a.VacancyID == uuid.Nil {
		return ErrInvalidInput
	}
	if a.Answer == nil {
		return ErrInvalidInput
	}
	return nil
}

func validateDepartment(d *models.Department) error {
	if d.Name == "" {
		return ErrInvalidInput
	}
	return nil
}

func validateLevel(l *models.Level) error {
	if l.Name == "" {
		return ErrInvalidInput
	}
	return nil
}

func validateLocation(l *models.Location) error {
	if l.Name == "" {
		return ErrInvalidInput
	}
	return nil
} 