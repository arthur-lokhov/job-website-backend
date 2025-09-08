// Package v1 contains the full implementation of the public API and admin panel
package v1

// PublicVacancyItem represents a vacancy item for public view.
type PublicVacancyItem struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Location   *LocationDTO   `json:"location"`
	Department *DepartmentDTO `json:"department"`
	Level      *LevelDTO      `json:"level"`
	LogoURL    string         `json:"logoUrl,omitempty"`
	Important  bool           `json:"important"`
	Priority   int            `json:"priority"`
}

// PublicVacancyDetails represents detailed information about a vacancy for public view.

type PublicVacancyDetails struct {
	PublicVacancyItem
	PhotoURL string `json:"photoUrl,omitempty"`
	Info     string `json:"info"`
}

// LocationDTO represents a location data transfer object.
type LocationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DepartmentDTO represents a department data transfer object.
type DepartmentDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// LevelDTO represents a level data transfer object.
type LevelDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ApplicationFormDTO represents an application form data transfer object.
type ApplicationFormDTO struct {
	FormSchema map[string]interface{} `json:"formSchema"`
}

// ApplicationSubmitRequest represents a request to submit an application.
type ApplicationSubmitRequest struct {
	VacancyID string                 `json:"vacancyId" validate:"required,uuid"`
	Answers   map[string]interface{} `json:"answers" validate:"required"`
}

// SuccessResponse represents a successful response.
type SuccessResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// AdminVacancyItem represents a vacancy item for admin view.
type AdminVacancyItem struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Location        *LocationDTO       `json:"location"`
	Department      *DepartmentDTO     `json:"department"`
	Level           *LevelDTO          `json:"level"`
	Important       bool               `json:"important"`
	Priority        int                `json:"priority"`
	Info            string             `json:"info"`
	ApplicationForm ApplicationFormDTO `json:"applicationForm"`
	IsActive        bool               `json:"isActive"`
	CreatedAt       string             `json:"createdAt"`
	UpdatedAt       string             `json:"updatedAt"`
}

// CreatedResponse represents a response for a created resource.
type CreatedResponse struct {
	Success   bool   `json:"success"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

// ApplicationItem represents an application item.
type ApplicationItem struct {
	ID          string                 `json:"id"`
	VacancyID   string                 `json:"vacancyId"`
	VacancyName string                 `json:"vacancyName"`
	Answers     map[string]interface{} `json:"answers"`
	Status      string                 `json:"status"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
	AnsweredAt  *string                `json:"answeredAt,omitempty"`
}

// ApplicationsList represents a list of application items.
type ApplicationsList []ApplicationItem

// ApplicationStatusUpdate represents a request to update an application's status.
type ApplicationStatusUpdate struct {
	Status string `json:"status" validate:"required,oneof=PENDING VIEWED INPROGRESS APPROVED REJECTED"`
}

// VacancyCreateRequest represents a request to create a vacancy.
type VacancyCreateRequest struct {
	Name            string              `json:"name" validate:"required,min=3,max=100"`
	LocationID      string              `json:"locationId" validate:"required,uuid"`
	DepartmentID    string              `json:"departmentId" validate:"required,uuid"`
	LevelID         string              `json:"levelId" validate:"required,uuid"`
	Info            string              `json:"info" validate:"required,min=10"`
	ApplicationForm *ApplicationFormDTO `json:"applicationForm" validate:"required"`
	Important       bool                `json:"important"`
	Priority        int                 `json:"priority" validate:"min=0"`
	IsActive        bool                `json:"isActive"`
}

// VacancyUpdateRequest represents a request to update a vacancy.
type VacancyUpdateRequest struct {
	Name            *string             `json:"name" validate:"omitempty,min=3,max=100"`
	LocationID      *string             `json:"locationId" validate:"omitempty,uuid"`
	DepartmentID    *string             `json:"departmentId" validate:"omitempty,uuid"`
	LevelID         *string             `json:"levelId" validate:"omitempty,uuid"`
	Info            *string             `json:"info" validate:"omitempty,min=10"`
	ApplicationForm *ApplicationFormDTO `json:"applicationForm"`
	Important       *bool               `json:"important"`
	Priority        *int                `json:"priority" validate:"omitempty,min=0"`
	IsActive        *bool               `json:"isActive"`
}
