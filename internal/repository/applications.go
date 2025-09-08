// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"job-website-backend/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewApplicationRepository(pool *pgxpool.Pool, log *slog.Logger) *ApplicationRepository {
	return &ApplicationRepository{Pool: pool, log: log}
}

func (r *ApplicationRepository) CreateApplication(ctx context.Context, vacancyID uuid.UUID, answers []byte) (string, error) {
	const op = "repository.ApplicationRepository.CreateApplication"
	log := r.log.With(slog.String("op", op), slog.String("vacancyID", vacancyID.String()))

	log.Info("creating application")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Insert("applications").
		Columns("vacancy_id", "answer", "status").
		Values(vacancyID, answers, "PENDING").
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return "", err
	}

	var appID string
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&appID); err != nil {
		log.Error("failed to create application", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully created application")
	return appID, nil
}

func (r *ApplicationRepository) GetApplicationsByAdmin(ctx context.Context, status, vacancyID, search string) ([]models.Application, error) {
	const op = "repository.ApplicationRepository.GetApplicationsByAdmin"
	log := r.log.With(slog.String("op", op), slog.String("status", status), slog.String("vacancyID", vacancyID), slog.String("search", search))

	log.Info("getting applications by admin")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query := psql.Select(
		"id",
		"vacancy_id",
		"answer",
		"status",
		"created_at",
		"updated_at",
		"answered_at",
	).From("applications")

	if status != "" {
		query = query.Where(squirrel.Eq{"status": status})
	}

	if vacancyID != "" {
		query = query.Where(squirrel.Eq{"vacancy_id": vacancyID})
	}

	if search != "" {
		// Assuming FIO is stored under a key like "fullName" or "name" in the JSON answer
		// Using jsonb_extract_path_text for case-insensitive search
		query = query.Where(squirrel.Like{"LOWER(answer->>\"fullName\")": "%" + strings.ToLower(search) + "%"})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Error("failed to query applications", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to query applications: %w", err)
	}
	defer rows.Close()

	var applications []models.Application
	for rows.Next() {
		var app models.Application
		if err := rows.Scan(
			&app.ID,
			&app.VacancyID,
			&app.Answer,
			&app.Status,
			&app.CreatedAt,
			&app.UpdatedAt,
			&app.AnsweredAt,
		); err != nil {
			log.Error("failed to scan application", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to scan application row: %w", err)
		}
		applications = append(applications, app)
	}

	if err := rows.Err(); err != nil {
		log.Error("error after iterating application rows", slog.String("error", err.Error()))
		return nil, fmt.Errorf("error after iterating application rows: %w", err)
	}

	log.Info("successfully got applications by admin")
	return applications, nil
}

func (r *ApplicationRepository) UpdateApplicationStatus(ctx context.Context, id string, status string) (string, error) {
	const op = "repository.ApplicationRepository.UpdateApplicationStatus"
	log := r.log.With(slog.String("op", op), slog.String("id", id), slog.String("status", status))

	log.Info("updating application status")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Update("applications").
		Set("status", status).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING \"id\"").
		ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return "", err
	}

	var appID string
	if err := r.Pool.QueryRow(ctx, sql, args...).Scan(&appID); err != nil {
		log.Error("failed to update application status", slog.String("error", err.Error()))
		return "", err
	}

	log.Info("successfully updated application status")
	return appID, nil
}

func (r *ApplicationRepository) DeleteApplication(ctx context.Context, id string) error {
	const op = "repository.ApplicationRepository.DeleteApplication"
	log := r.log.With(slog.String("op", op), slog.String("id", id))

	log.Info("deleting application")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := psql.Delete("applications").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		log.Error("failed to delete application", slog.String("error", err.Error()))
		return err
	}

	log.Info("successfully deleted application")
	return nil
}
