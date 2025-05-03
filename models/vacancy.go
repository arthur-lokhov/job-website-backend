package model

// TODO: Добавить UUID, FK, JSONB, time

type Vacancy struct {
	ID              string // need UUID
	Name            string
	Department      any // need FK
	Level           any // need FK
	Location        any // need FK
	Info            string
	ApplicationForm any // need JSONB
	IsActive        bool
	Important       bool
	Priority        int
	CreatedAt       any // need time
	UpdatedAt       any // need time
}
