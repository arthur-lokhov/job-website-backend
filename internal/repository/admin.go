// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"job-website-backend/internal/dto"
	"job-website-backend/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewAdminRepository(pool *pgxpool.Pool, log *slog.Logger) *AdminRepository {
	return &AdminRepository{Pool: pool, log: log}
}

func (r *AdminRepository) GetVacancies(ctx context.Context, isActive bool) ([]models.Vacancy, error) {
	const op = "repository.AdminRepository.GetVacancies"
	log := r.log.With(slog.String("op", op), slog.Bool("isActive", isActive))

	log.Info("getting vacancies")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := psql.Select(
		"v.id",
		"v.name",
		"l.id as location_id",
		"l.name as location_name",
		"d.id as department_id",
		"lv.id as level_id",
		"lv.name as level_name",
		"v.important",
		"v.priority",
		"v.info",
		"v.application_form",
		"v.is_active",
		"v.created_at",
		"v.updated_at",
	).From("vacancies v").
		LeftJoin("locations l ON v.location_id = l.id").
		LeftJoin("departments d ON v.department_id = d.id").
		LeftJoin("levels lv ON v.level_id = lv.id")

	if isActive {
		query = query.Where(squirrel.Eq{"v.is_active": true})
	} else {
		query = query.Where(squirrel.Eq{"v.is_active": false})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Error("failed to query vacancies", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var vacancies []models.Vacancy
	for rows.Next() {
		var v models.Vacancy
		var appFormRaw []byte
		v.Location = models.Location{}
		v.Level = models.Level{}
		if err := rows.Scan(
			&v.ID,
			&v.Name,
			&v.Location.ID,
			&v.Location.Name,
			&v.DepartmentID,
			&v.Level.ID,
			&v.Level.Name,
			&v.Important,
			&v.Priority,
			&v.Info,
			&appFormRaw,
			&v.IsActive,
			&v.CreatedAt,
			&v.UpdatedAt,
		); err != nil {
			log.Error("failed to scan vacancy", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	log.Info("successfully got vacancies")
	return vacancies, nil
}

func (r *AdminRepository) CreateVacancy(ctx context.Context, req dto.VacancyCreateDTO) (string, time.Time, error) {
	const op = "repository.AdminRepository.CreateVacancy"
	log := r.log.With(slog.String("op", op))

	log.Info("creating vacancy")

	form, err := json.Marshal(req.ApplicationForm)
	if err != nil {
		log.Error("failed to marshal application form", slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Insert("vacancies").
		Columns("name", "location_id", "department_id", "level_id", "important", "priority", "info", "application_form", "is_active").
		Values(req.Name, req.LocationID, req.DepartmentID, req.LevelID, req.Important, req.Priority, req.Info, form, req.IsActive).
		Suffix("RETURNING \"id\", \"created_at\"").
		ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	var vacancyID string
	var createdAt time.Time
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&vacancyID, &createdAt); err != nil {
		log.Error("failed to create vacancy", slog.String("error", err.Error()))
		return "", time.Time{}, err
	}

	log.Info("successfully created vacancy")
	return vacancyID, createdAt, nil
}

func (r *AdminRepository) UpdateVacancy(ctx context.Context, id string, req dto.VacancyUpdateDTO) (string, error) {
	const op = "repository.AdminRepository.UpdateVacancy"
	log := r.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("updating vacancy")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	updateBuilder := psql.Update("vacancies")

	if req.Name != "" {
		updateBuilder = updateBuilder.Set("name", req.Name)
	}
	if req.LocationID != "" {
		updateBuilder = updateBuilder.Set("location_id", req.LocationID)
	}
	if req.DepartmentID != "" {
		updateBuilder = updateBuilder.Set("department_id", req.DepartmentID)
	}
	if req.LevelID != "" {
		updateBuilder = updateBuilder.Set("level_id", req.LevelID)
	}
	if req.Important != nil {
		updateBuilder = updateBuilder.Set("important", *req.Important)
	}
	if req.Priority != nil {
		updateBuilder = updateBuilder.Set("priority", *req.Priority)
	}
	if req.Info != "" {
		updateBuilder = updateBuilder.Set("info", req.Info)
	}
	if req.ApplicationForm != nil {
		form, err := json.Marshal(req.ApplicationForm)
		if err != nil {
			log.Error("failed to marshal application form", slog.String("error", err.Error()))
			return "", err
		}
		updateBuilder = updateBuilder.Set("application_form", form)
	}
	if req.IsActive != nil {
		updateBuilder = updateBuilder.Set("is_active", *req.IsActive)
	}

	sql, args, err := updateBuilder.Where(squirrel.Eq{"id": id}).Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return "", err
	}

	var vacancyID string
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&vacancyID); err != nil {
		log.Error("failed to update vacancy", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully updated vacancy")
	return vacancyID, nil
}

func (r *AdminRepository) DeleteVacancy(ctx context.Context, id string) error {
	const op = "repository.AdminRepository.DeleteVacancy"
	log := r.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("deleting vacancy")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Delete("vacancies").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		log.Error("failed to delete vacancy", slog.String("error", err.Error()))
		return err
	}

	log.Info("successfully deleted vacancy")
	return nil
}
