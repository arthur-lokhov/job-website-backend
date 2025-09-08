// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"log/slog"

	"job-website-backend/internal/dto"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DepartmentRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewDepartmentRepository(pool *pgxpool.Pool, log *slog.Logger) *DepartmentRepository {
	return &DepartmentRepository{Pool: pool, log: log}
}

func (r *DepartmentRepository) UpsertDepartments(ctx context.Context, departments []dto.DepartmentDTO) error {
	const op = "repository.DepartmentRepository.UpsertDepartments"
	log := r.log.With(slog.String("op", op))

	log.Info("upserting departments")

	if len(departments) == 0 {
		log.Info("no departments to upsert")
		return nil
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	insertBuilder := psql.Insert("departments").Columns("id")

	for _, dept := range departments {
		insertBuilder = insertBuilder.Values(dept.ID)
	}

	sql, args, err := insertBuilder.Suffix("ON CONFLICT (id) DO NOTHING").ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return err
	}

	if _, err := r.Pool.Exec(ctx, sql, args...); err != nil {
		log.Error("failed to upsert departments", slog.String("error", err.Error()))
		return err
	}

	log.Info("successfully upserted departments")
	return nil
}
