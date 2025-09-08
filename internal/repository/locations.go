// Package repository contains the database access logic for the application.
package repository

import (
	"context"
	"log/slog"

	"job-website-backend/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	Pool *pgxpool.Pool
	log  *slog.Logger
}

func NewLocationRepository(pool *pgxpool.Pool, log *slog.Logger) *LocationRepository {
	return &LocationRepository{Pool: pool, log: log}
}

func (r *LocationRepository) GetLocations(ctx context.Context) ([]models.Location, error) {
	const op = "repository.LocationRepository.GetLocations"
	log := r.log.With(slog.String("op", op))

	log.Info("getting locations")

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := psql.Select("id", "name").From("locations").ToSql()
	if err != nil {
		log.Error("failed to build sql", slog.String("error", err.Error()))
		return nil, err
	}

	rows, err := r.Pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to query locations", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var locations []models.Location
	for rows.Next() {
		var l models.Location
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			log.Error("failed to scan location", slog.String("error", err.Error()))
			return nil, err
		}
		locations = append(locations, l)
	}

	log.Info("successfully got locations")
	return locations, nil
}
