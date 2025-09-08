// Package models contains the data models for the application.
package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID              uuid.UUID
	Name            string
	DepartmentID    uuid.UUID
	DepartmentName  string // Populated from cache
	LevelID         uuid.UUID
	Level           Level
	LocationID      uuid.UUID
	Location        Location
	Info            string // markdown
	ApplicationForm json.RawMessage
	IsActive        bool
	Important       bool
	Priority        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Application struct {
	ID         uuid.UUID
	VacancyID  uuid.UUID
	Vacancy    Vacancy
	Answer     json.RawMessage
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AnsweredAt *time.Time
}

type Department struct {
	ID uuid.UUID
}

type Level struct {
	ID        uuid.UUID
	Name      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Location struct {
	ID        uuid.UUID
	Name      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
