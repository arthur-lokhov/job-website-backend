// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LevelRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewLevelRepository(pool *pgxpool.Pool, log *slog.Logger) *LevelRepository {
	return &LevelRepository{Pool: pool, log: log}
}

func (r *LevelRepository) GetLevels(ctx context.Context) ([]models.Level, error) {
	const op = "repository.LevelRepository.GetLevels"
	log := r.log.With(slog.String("op", op))

	log.Info("getting levels")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.Select("id", "name").From("levels").ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to query levels", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var levels []models.Level
	for rows.Next() {
		var l models.Level
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			log.Error("failed to scan level", slog.String("error", err.Error()))
			return nil, err
		}
		levels = append(levels, l)
	}

	log.Info("successfully got levels")
	return levels, nil
}
