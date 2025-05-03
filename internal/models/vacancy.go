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

	DepartmentID uuid.UUID
	Department   Department `gorm:"foreignKey:DepartmentID"`

	LevelID uuid.UUID
	Level   Level `gorm:"foreignKey:LevelID"`

	LocationID uuid.UUID
	Location   Location `gorm:"foreignKey:LocationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
