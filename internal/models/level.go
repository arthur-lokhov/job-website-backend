package models

import (
	"time"

	"github.com/google/uuid"
)

type Level struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
