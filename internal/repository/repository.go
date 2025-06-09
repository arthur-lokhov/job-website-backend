package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/models"
)

type Repository interface {
	// Vacancy operations
	CreateVacancy(ctx context.Context, vacancy *models.Vacancy) error
	GetVacancy(ctx context.Context, id uuid.UUID) (*models.Vacancy, error)
	UpdateVacancy(ctx context.Context, vacancy *models.Vacancy) error
	DeleteVacancy(ctx context.Context, id uuid.UUID) error
	ListVacancies(ctx context.Context, isActive *bool, important *bool) ([]*models.Vacancy, error)
	SearchVacancies(ctx context.Context, query string) ([]*models.Vacancy, error)

	// Application operations
	CreateApplication(ctx context.Context, application *models.Application) error
	GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error)
	UpdateApplication(ctx context.Context, application *models.Application) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error
	ListApplications(ctx context.Context, status *string, vacancyID *uuid.UUID) ([]*models.Application, error)

	// Department operations
	CreateDepartment(ctx context.Context, department *models.Department) error
	GetDepartment(ctx context.Context, id uuid.UUID) (*models.Department, error)
	UpdateDepartment(ctx context.Context, department *models.Department) error
	DeleteDepartment(ctx context.Context, id uuid.UUID) error
	ListDepartments(ctx context.Context) ([]*models.Department, error)

	// Level operations
	CreateLevel(ctx context.Context, level *models.Level) error
	GetLevel(ctx context.Context, id uuid.UUID) (*models.Level, error)
	UpdateLevel(ctx context.Context, level *models.Level) error
	DeleteLevel(ctx context.Context, id uuid.UUID) error
	ListLevels(ctx context.Context) ([]*models.Level, error)

	// Location operations
	CreateLocation(ctx context.Context, location *models.Location) error
	GetLocation(ctx context.Context, id uuid.UUID) (*models.Location, error)
	UpdateLocation(ctx context.Context, location *models.Location) error
	DeleteLocation(ctx context.Context, id uuid.UUID) error
	ListLocations(ctx context.Context) ([]*models.Location, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

// Helper functions
func (r *repository) scanVacancy(row pgx.Row) (*models.Vacancy, error) {
	var v models.Vacancy
	var form []byte
	err := row.Scan(
		&v.ID,
		&v.Name,
		&v.DepartmentID,
		&v.LevelID,
		&v.LocationID,
		&v.Info,
		&form,
		&v.IsActive,
		&v.Important,
		&v.Priority,
		&v.CreatedAt,
		&v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	v.ApplicationForm = form
	return &v, nil
}

func (r *repository) scanApplication(row pgx.Row) (*models.Application, error) {
	var a models.Application
	var answer []byte
	err := row.Scan(
		&a.ID,
		&a.VacancyID,
		&answer,
		&a.Status,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.AnsweredAt,
	)
	if err != nil {
		return nil, err
	}
	a.Answer = answer
	return &a, nil
}

// Vacancy operations
func (r *repository) CreateVacancy(ctx context.Context, vacancy *models.Vacancy) error {
	query := `
		INSERT INTO vacancies (
			name, department_id, level_id, location_id, info, 
			application_form, is_active, important, priority
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		vacancy.Name,
		vacancy.DepartmentID,
		vacancy.LevelID,
		vacancy.LocationID,
		vacancy.Info,
		vacancy.ApplicationForm,
		vacancy.IsActive,
		vacancy.Important,
		vacancy.Priority,
	).Scan(&vacancy.ID, &vacancy.CreatedAt, &vacancy.UpdatedAt)
}

func (r *repository) GetVacancy(ctx context.Context, id uuid.UUID) (*models.Vacancy, error) {
	query := `
		SELECT id, name, department_id, level_id, location_id, info,
			application_form, is_active, important, priority,
			created_at, updated_at
		FROM vacancies
		WHERE id = $1`

	return r.scanVacancy(r.pool.QueryRow(ctx, query, id))
}

func (r *repository) UpdateVacancy(ctx context.Context, vacancy *models.Vacancy) error {
	query := `
		UPDATE vacancies
		SET name = $1, department_id = $2, level_id = $3, location_id = $4,
			info = $5, application_form = $6, is_active = $7,
			important = $8, priority = $9
		WHERE id = $10
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		vacancy.Name,
		vacancy.DepartmentID,
		vacancy.LevelID,
		vacancy.LocationID,
		vacancy.Info,
		vacancy.ApplicationForm,
		vacancy.IsActive,
		vacancy.Important,
		vacancy.Priority,
		vacancy.ID,
	).Scan(&vacancy.UpdatedAt)
}

func (r *repository) DeleteVacancy(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM vacancies WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListVacancies(ctx context.Context, isActive *bool, important *bool) ([]*models.Vacancy, error) {
	query := `
		SELECT id, name, department_id, level_id, location_id, info,
			application_form, is_active, important, priority,
			created_at, updated_at
		FROM vacancies
		WHERE ($1::boolean IS NULL OR is_active = $1)
		AND ($2::boolean IS NULL OR important = $2)
		ORDER BY priority DESC, created_at DESC`

	rows, err := r.pool.Query(ctx, query, isActive, important)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []*models.Vacancy
	for rows.Next() {
		vacancy, err := r.scanVacancy(rows)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, vacancy)
	}
	return vacancies, rows.Err()
}

func (r *repository) SearchVacancies(ctx context.Context, query string) ([]*models.Vacancy, error) {
	sqlQuery := `
		SELECT id, name, department_id, level_id, location_id, info,
			application_form, is_active, important, priority,
			created_at, updated_at
		FROM vacancies
		WHERE is_active = true
		AND (name ILIKE $1 OR info ILIKE $1)
		ORDER BY priority DESC, created_at DESC`

	rows, err := r.pool.Query(ctx, sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []*models.Vacancy
	for rows.Next() {
		vacancy, err := r.scanVacancy(rows)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, vacancy)
	}
	return vacancies, rows.Err()
}

// Application operations
func (r *repository) CreateApplication(ctx context.Context, application *models.Application) error {
	query := `
		INSERT INTO applications (
			vacancy_id, answer, status
		) VALUES (
			$1, $2, $3
		) RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		application.VacancyID,
		application.Answer,
		application.Status,
	).Scan(&application.ID, &application.CreatedAt, &application.UpdatedAt)
}

func (r *repository) GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	query := `
		SELECT id, vacancy_id, answer, status,
			created_at, updated_at, answered_at
		FROM applications
		WHERE id = $1`

	return r.scanApplication(r.pool.QueryRow(ctx, query, id))
}

func (r *repository) UpdateApplication(ctx context.Context, application *models.Application) error {
	query := `
		UPDATE applications
		SET answer = $1, status = $2, answered_at = $3
		WHERE id = $4
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		application.Answer,
		application.Status,
		application.AnsweredAt,
		application.ID,
	).Scan(&application.UpdatedAt)
}

func (r *repository) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM applications WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListApplications(ctx context.Context, status *string, vacancyID *uuid.UUID) ([]*models.Application, error) {
	query := `
		SELECT id, vacancy_id, answer, status,
			created_at, updated_at, answered_at
		FROM applications
		WHERE ($1::text IS NULL OR status = $1)
		AND ($2::uuid IS NULL OR vacancy_id = $2)
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, status, vacancyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []*models.Application
	for rows.Next() {
		application, err := r.scanApplication(rows)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}
	return applications, rows.Err()
}

// Department operations
func (r *repository) CreateDepartment(ctx context.Context, department *models.Department) error {
	query := `
		INSERT INTO departments (name)
		VALUES ($1)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query, department.Name).
		Scan(&department.ID, &department.CreatedAt, &department.UpdatedAt)
}

func (r *repository) GetDepartment(ctx context.Context, id uuid.UUID) (*models.Department, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM departments
		WHERE id = $1`

	var d models.Department
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repository) UpdateDepartment(ctx context.Context, department *models.Department) error {
	query := `
		UPDATE departments
		SET name = $1
		WHERE id = $2
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query, department.Name, department.ID).
		Scan(&department.UpdatedAt)
}

func (r *repository) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM departments WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListDepartments(ctx context.Context) ([]*models.Department, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM departments
		ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &d)
	}
	return departments, rows.Err()
}

// Level operations
func (r *repository) CreateLevel(ctx context.Context, level *models.Level) error {
	query := `
		INSERT INTO levels (name, priority)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query, level.Name, level.Priority).
		Scan(&level.ID, &level.CreatedAt, &level.UpdatedAt)
}

func (r *repository) GetLevel(ctx context.Context, id uuid.UUID) (*models.Level, error) {
	query := `
		SELECT id, name, priority, created_at, updated_at
		FROM levels
		WHERE id = $1`

	var l models.Level
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&l.ID, &l.Name, &l.Priority, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *repository) UpdateLevel(ctx context.Context, level *models.Level) error {
	query := `
		UPDATE levels
		SET name = $1, priority = $2
		WHERE id = $3
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query, level.Name, level.Priority, level.ID).
		Scan(&level.UpdatedAt)
}

func (r *repository) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM levels WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListLevels(ctx context.Context) ([]*models.Level, error) {
	query := `
		SELECT id, name, priority, created_at, updated_at
		FROM levels
		ORDER BY priority DESC, name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var levels []*models.Level
	for rows.Next() {
		var l models.Level
		err := rows.Scan(&l.ID, &l.Name, &l.Priority, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		levels = append(levels, &l)
	}
	return levels, rows.Err()
}

// Location operations
func (r *repository) CreateLocation(ctx context.Context, location *models.Location) error {
	query := `
		INSERT INTO locations (name, priority)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query, location.Name, location.Priority).
		Scan(&location.ID, &location.CreatedAt, &location.UpdatedAt)
}

func (r *repository) GetLocation(ctx context.Context, id uuid.UUID) (*models.Location, error) {
	query := `
		SELECT id, name, priority, created_at, updated_at
		FROM locations
		WHERE id = $1`

	var l models.Location
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&l.ID, &l.Name, &l.Priority, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *repository) UpdateLocation(ctx context.Context, location *models.Location) error {
	query := `
		UPDATE locations
		SET name = $1, priority = $2
		WHERE id = $3
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query, location.Name, location.Priority, location.ID).
		Scan(&location.UpdatedAt)
}

func (r *repository) DeleteLocation(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM locations WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *repository) ListLocations(ctx context.Context) ([]*models.Location, error) {
	query := `
		SELECT id, name, priority, created_at, updated_at
		FROM locations
		ORDER BY priority DESC, name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []*models.Location
	for rows.Next() {
		var l models.Location
		err := rows.Scan(&l.ID, &l.Name, &l.Priority, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &l)
	}
	return locations, rows.Err()
} 