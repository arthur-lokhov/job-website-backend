// Package dto contains data transfer objects for the application.
package dto

import "github.com/google/uuid"

// ApplicationFormDTO defines the structure for the application form schema.
// It is used in both VacancyCreateDTO and VacancyUpdateDTO.
type ApplicationFormDTO struct {
	FormSchema map[string]any `json:"formSchema"`
}

// VacancyCreateDTO is used to transfer data from the controller to the service for creating a vacancy.
// It is a copy of the v1.VacancyCreateRequest struct.
type VacancyCreateDTO struct {
	Name            string              `json:"name"`
	LocationID      string              `json:"locationId"`
	DepartmentID    string              `json:"departmentId"`
	LevelID         string              `json:"levelId"`
	Important       bool                `json:"important"`
	Priority        int                 `json:"priority"`
	Info            string              `json:"info"`
	ApplicationForm *ApplicationFormDTO `json:"applicationForm"`
	IsActive        bool                `json:"isActive"`
}

// VacancyUpdateDTO is used to transfer data from the controller to the service for updating a vacancy.
// It is a copy of the v1.VacancyUpdateRequest struct.
type VacancyUpdateDTO struct {
	Name            string              `json:"name,omitempty"`
	LocationID      string              `json:"locationId,omitempty"`
	DepartmentID    string              `json:"departmentId,omitempty"`
	LevelID         string              `json:"levelId,omitempty"`
	Important       *bool               `json:"important,omitempty"`
	Priority        *int                `json:"priority,omitempty"`
	Info            string              `json:"info,omitempty"`
	ApplicationForm *ApplicationFormDTO `json:"applicationForm,omitempty"`
	IsActive        *bool               `json:"isActive,omitempty"`
}

// DepartmentDTO defines the structure for department data transfer.
type DepartmentDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}
