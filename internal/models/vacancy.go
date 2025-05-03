package models

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string
	Description string
	CompanyName string
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
