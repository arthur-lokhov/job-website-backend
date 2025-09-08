// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VacancyRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewVacancyRepository(pool *pgxpool.Pool, log *slog.Logger) *VacancyRepository {
	return &VacancyRepository{Pool: pool, log: log}
}

func (r *VacancyRepository) GetVacancies(ctx context.Context, important bool, sort string, search string) ([]models.Vacancy, error) {
	const op = "repository.VacancyRepository.GetVacancies"
	log := r.log.With(slog.String("op", op), slog.Bool("important", important), slog.String("sort", sort), slog.String("search", search))

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
	).From("vacancies v").
		LeftJoin("locations l ON v.location_id = l.id").
		LeftJoin("departments d ON v.department_id = d.id").
		LeftJoin("levels lv ON v.level_id = lv.id").
		Where(squirrel.Eq{"v.is_active": true})

	if important {
		query = query.Where(squirrel.Eq{"v.important": true})
	}

	if search != "" {
		query = query.Where(squirrel.Like{"v.name": "%" + search + "%"})
	}

	if sort == "asc" {
		query = query.OrderBy("v.priority ASC")
	} else {
		query = query.OrderBy("v.priority DESC")
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
		); err != nil {
			log.Error("failed to scan vacancy", slog.String("error", err.Error()))
			return nil, err
		}
		vacancies = append(vacancies, v)
	}

	log.Info("successfully got vacancies")
	return vacancies, nil
}

func (r *VacancyRepository) GetVacancyByID(ctx context.Context, id string) (*models.Vacancy, error) {
	const op = "repository.VacancyRepository.GetVacancyByID"
	log := r.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("getting vacancy by id")

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
	).From("vacancies v").
		LeftJoin("locations l ON v.location_id = l.id").
		LeftJoin("departments d ON v.department_id = d.id").
		LeftJoin("levels lv ON v.level_id = lv.id").
		Where(squirrel.Eq{"v.id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, err
	}

	var v models.Vacancy
	v.Location = models.Location{}
	v.Level = models.Level{}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
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
	)
	if err != nil {
		log.Error("failed to query vacancy by id", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got vacancy by id")
	return &v, nil
}

func (r *VacancyRepository) GetVacancyForm(ctx context.Context, id string) ([]byte, error) {
	const op = "repository.VacancyRepository.GetVacancyForm"
	log := r.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("getting vacancy form")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Select("application_form").From("vacancies").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, err
	}

	var applicationForm []byte
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&applicationForm)
	if err != nil {
		log.Error("failed to query vacancy form", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("successfully got vacancy form")
	return applicationForm, nil
}
