package models

import (
	"github.com/google/uuid"
	"time"
	"gorm.io/datatypes"
)

type Vacancy struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Name            string
	DepartmentID    uuid.UUID
	Department      Department
	LevelID         uuid.UUID
	Level           Level
	LocationID      uuid.UUID
	Location        Location
	Info            string         // markdown
	ApplicationForm datatypes.JSON
	IsActive        bool
	Important       bool
	Priority        int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Application struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey"`
	VacancyID  uuid.UUID
	Vacancy    Vacancy
	Answer     datatypes.JSON
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AnsweredAt *time.Time
}

type Department struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

type Level struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Location struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthService struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Slug      string
	LogoUrl   string
	PhotoUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
} 