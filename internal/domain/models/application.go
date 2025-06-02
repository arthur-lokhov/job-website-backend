package models

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`

	VacancyID uuid.UUID
	Vacancy   Vacancy `gorm:"foreignKey:VacancyID"`

	Answer     string `gorm:"type:jsonb"`
	Status     string `gorm:"type:application_status"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AnsweredAt time.Time
}
