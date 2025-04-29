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
	isActive        bool
	Important       bool
	priority        int
	createdAt       any // need time
	updatedAt       any // need time
}
