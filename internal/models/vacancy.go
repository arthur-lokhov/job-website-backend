package models

import (
	"time"

	"github.com/google/uuid"
)

type Vacancy struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Info      string
	IsActive  bool
	Important bool
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time

	// Department string
	// Level string
	// Location  string
	// ApplicationForm json
}
