package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Vacancy represents a job vacancy
type Vacancy struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	Name            string      `json:"name" db:"name"`
	DepartmentID    uuid.UUID   `json:"departmentId" db:"department_id"`
	LevelID         uuid.UUID   `json:"levelId" db:"level_id"`
	LocationID      uuid.UUID   `json:"locationId" db:"location_id"`
	Info            string      `json:"info" db:"info"`
	ApplicationForm interface{} `json:"applicationForm" db:"application_form"`
	IsActive        bool        `json:"isActive" db:"is_active"`
	Important       bool        `json:"important" db:"important"`
	Priority        int         `json:"priority" db:"priority"`
	CreatedAt       time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time   `json:"updatedAt" db:"updated_at"`
}

// Department represents a company department
type Department struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Application represents a job application
type Application struct {
	ID         uuid.UUID   `json:"id" db:"id"`
	VacancyID  uuid.UUID   `json:"vacancyId" db:"vacancy_id"`
	Answer     interface{} `json:"answer" db:"answer"`
	Status     string      `json:"status" db:"status"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time   `json:"updatedAt" db:"updated_at"`
	AnsweredAt *time.Time  `json:"answeredAt,omitempty" db:"answered_at"`
}

// Error represents an API error response
type Error struct {
	Error   string `json:"error" example:"Bad Request"`
	Message string `json:"message" example:"Invalid request parameters"`
}

// CreateVacancyRequest represents the request body for creating a vacancy
type CreateVacancyRequest struct {
	Name            string      `json:"name" example:"Senior Software Engineer"`
	DepartmentID    uuid.UUID   `json:"departmentId" example:"123e4567-e89b-12d3-a456-426614174000"`
	LevelID         uuid.UUID   `json:"levelId" example:"123e4567-e89b-12d3-a456-426614174000"`
	LocationID      uuid.UUID   `json:"locationId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Info            string      `json:"info" example:"We are looking for an experienced software engineer..."`
	ApplicationForm interface{} `json:"applicationForm" example:"{\"type\":\"object\",\"properties\":{\"name\":{\"type\":\"string\"}}}"`
	IsActive        bool        `json:"isActive" example:"true"`
	Important       bool        `json:"important" example:"false"`
	Priority        int         `json:"priority" example:"1"`
}

// UpdateVacancyRequest represents the request body for updating a vacancy
type UpdateVacancyRequest struct {
	Name            string      `json:"name" example:"Senior Software Engineer"`
	DepartmentID    uuid.UUID   `json:"departmentId" example:"123e4567-e89b-12d3-a456-426614174000"`
	LevelID         uuid.UUID   `json:"levelId" example:"123e4567-e89b-12d3-a456-426614174000"`
	LocationID      uuid.UUID   `json:"locationId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Info            string      `json:"info" example:"We are looking for an experienced software engineer..."`
	ApplicationForm interface{} `json:"applicationForm" example:"{\"type\":\"object\",\"properties\":{\"name\":{\"type\":\"string\"}}}"`
	IsActive        bool        `json:"isActive" example:"true"`
	Important       bool        `json:"important" example:"false"`
	Priority        int         `json:"priority" example:"1"`
}

// SubmitApplicationRequest represents the request body for submitting an application
type SubmitApplicationRequest struct {
	VacancyID uuid.UUID   `json:"vacancyId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Answer    interface{} `json:"answer" example:"{\"name\":\"John Doe\",\"email\":\"john@example.com\"}"`
}

// UpdateApplicationStatusRequest represents the request body for updating an application status
type UpdateApplicationStatusRequest struct {
	Status string `json:"status" example:"APPROVED" enums:"PENDING,APPROVED,REJECTED"`
}

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool      `json:"success" example:"true"`
	ID      uuid.UUID `json:"id,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// CreateVacancyResponse represents the response for creating a vacancy
type CreateVacancyResponse struct {
	Success   bool      `json:"success" example:"true"`
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedAt time.Time `json:"createdAt" example:"2024-03-20T10:00:00Z"`
}

// UpdateVacancyResponse represents the response for updating a vacancy
type UpdateVacancyResponse struct {
	Success bool      `json:"success" example:"true"`
	ID      uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// SubmitApplicationResponse represents the response for submitting an application
type SubmitApplicationResponse struct {
	Success bool      `json:"success" example:"true"`
	ID      uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// UpdateApplicationStatusResponse represents the response for updating an application status
type UpdateApplicationStatusResponse struct {
	Success bool      `json:"success" example:"true"`
	ID      uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type Level struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Priority  int       `json:"priority" db:"priority"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Location struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Priority  int       `json:"priority" db:"priority"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type AuthService struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	LogoURL   string    `json:"logoUrl" db:"logo_url"`
	PhotoURL  string    `json:"photoUrl" db:"photo_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Form components
type FormComponent struct {
	Type    string          `json:"type"`
	Config  json.RawMessage `json:"config"`
	Answer  interface{}     `json:"answer,omitempty"`
}

// Form component configs
type InputConfig struct {
	Required bool   `json:"required"`
	Label    string `json:"label"`
	Caption  string `json:"caption,omitempty"`
}

type TextareaConfig struct {
	Required bool   `json:"required"`
	Label    string `json:"label"`
}

type HeaderConfig struct {
	Label string `json:"label"`
}

type SelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type SelectConfig struct {
	Required bool           `json:"required"`
	Label    string         `json:"label"`
	Options  []SelectOption `json:"options"`
}

type RadioConfig struct {
	Required     bool           `json:"required"`
	DefaultValue string         `json:"defaultValue,omitempty"`
	Options      []SelectOption `json:"options"`
}

type CheckboxConfig struct {
	Required     bool           `json:"required"`
	DefaultValue string         `json:"defaultValue,omitempty"`
	Options      []SelectOption `json:"options"`
}