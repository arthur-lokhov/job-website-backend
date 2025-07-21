package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"job-website-backend/internal/models"
	"strings"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

var DB *gorm.DB // должен быть инициализирован при старте приложения

// --- DTOs для публичных ручек ---
type PublicVacancyItem struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Location  *LocationDTO   `json:"location"`
	Department *DepartmentDTO `json:"department"`
	Level     *LevelDTO      `json:"level"`
	LogoUrl   string         `json:"logoUrl,omitempty"`
	Important bool           `json:"important"`
	Priority  int            `json:"priority"`
}

type PublicVacancyDetails struct {
	PublicVacancyItem
	PhotoUrl string `json:"photoUrl,omitempty"`
	Info     string `json:"info"`
}

type LocationDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DepartmentDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LevelDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationFormDTO struct {
	FormSchema map[string]interface{} `json:"formSchema"`
}

type ApplicationSubmitRequest struct {
	VacancyID string                 `json:"vacancyId"`
	Answers   map[string]interface{} `json:"answers"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

// --- DTOs для админских ручек ---
type AdminVacancyItem struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Location        *LocationDTO   `json:"location"`
	Department      *DepartmentDTO `json:"department"`
	Level           *LevelDTO      `json:"level"`
	Important       bool           `json:"important"`
	Priority        int            `json:"priority"`
	Info            string         `json:"info"`
	ApplicationForm ApplicationFormDTO `json:"applicationForm"`
	IsActive        bool           `json:"isActive"`
	CreatedAt       string         `json:"createdAt"`
	UpdatedAt       string         `json:"updatedAt"`
}

type CreatedResponse struct {
	Success   bool   `json:"success"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type ApplicationItem struct {
	ID         string                 `json:"id"`
	VacancyID  string                 `json:"vacancyId"`
	VacancyName string                `json:"vacancyName"`
	Answers    map[string]interface{} `json:"answers"`
	Status     string                 `json:"status"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
	AnsweredAt *string                `json:"answeredAt,omitempty"`
}

type ApplicationsList []ApplicationItem

type ApplicationStatusUpdate struct {
	Status string `json:"status"`
}

// --- Публичные ручки ---
// GET /vacancies
func GetVacancies(w http.ResponseWriter, r *http.Request) {
	var vacancies []models.Vacancy
	query := DB.Preload("Location").Preload("Department").Preload("Level").Where("is_active = ?", true)

	// Фильтр important
	if imp := r.URL.Query().Get("important"); imp != "" {
		if strings.ToLower(imp) == "true" {
			query = query.Where("important = ?", true)
		} else if strings.ToLower(imp) == "false" {
			query = query.Where("important = ?", false)
		}
	}
	// Сортировка по приоритету
	sort := r.URL.Query().Get("sort")
	if sort == "asc" {
		query = query.Order("priority ASC")
	} else {
		query = query.Order("priority DESC")
	}
	if err := query.Find(&vacancies).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var resp []PublicVacancyItem
	for _, v := range vacancies {
		resp = append(resp, PublicVacancyItem{
			ID:        v.ID.String(),
			Name:      v.Name,
			Location:  &LocationDTO{ID: v.LocationID.String(), Name: v.Location.Name},
			Department: &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.Department.Name},
			Level:     &LevelDTO{ID: v.LevelID.String(), Name: v.Level.Name},
			LogoUrl:   "", // можно добавить если есть
			Important: v.Important,
			Priority:  v.Priority,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /vacancies/{id}
func GetVacancyByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var v models.Vacancy
	if err := DB.Preload("Location").Preload("Department").Preload("Level").First(&v, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Resource not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	resp := PublicVacancyDetails{
		PublicVacancyItem: PublicVacancyItem{
			ID:        v.ID.String(),
			Name:      v.Name,
			Location:  &LocationDTO{ID: v.LocationID.String(), Name: v.Location.Name},
			Department: &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.Department.Name},
			Level:     &LevelDTO{ID: v.LevelID.String(), Name: v.Level.Name},
			LogoUrl:   "",
			Important: v.Important,
			Priority:  v.Priority,
		},
		PhotoUrl: "",
		Info:     v.Info,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /vacancies/{id}/form
func GetVacancyForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var v models.Vacancy
	if err := DB.First(&v, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Resource not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var form map[string]interface{}
	_ = json.Unmarshal(v.ApplicationForm, &form)
	resp := ApplicationFormDTO{FormSchema: form}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /applications
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	var req ApplicationSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	if req.VacancyID == "" || req.Answers == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	vacID, err := uuid.Parse(req.VacancyID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid vacancyId"})
		return
	}
	ans, err := json.Marshal(req.Answers)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid answers format"})
		return
	}
	app := models.Application{
		VacancyID: vacID,
		Answer:   datatypes.JSON(ans),
		Status:   "PENDING",
	}
	if err := DB.Create(&app).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{Success: true, ID: app.ID.String()})
}

// Приватные
// --- Админские ручки ---
// GET /admin/vacancies
func AdminGetVacancies(w http.ResponseWriter, r *http.Request) {
	var vacancies []models.Vacancy
	query := DB.Preload("Location").Preload("Department").Preload("Level")
	if isActive := r.URL.Query().Get("isActive"); isActive != "" {
		if isActive == "true" {
			query = query.Where("is_active = ?", true)
		} else if isActive == "false" {
			query = query.Where("is_active = ?", false)
		}
	}
	if err := query.Find(&vacancies).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var resp []AdminVacancyItem
	for _, v := range vacancies {
		var form map[string]interface{}
		_ = json.Unmarshal(v.ApplicationForm, &form)
		resp = append(resp, AdminVacancyItem{
			ID:         v.ID.String(),
			Name:       v.Name,
			Location:   &LocationDTO{ID: v.LocationID.String(), Name: v.Location.Name},
			Department: &DepartmentDTO{ID: v.DepartmentID.String(), Name: v.Department.Name},
			Level:      &LevelDTO{ID: v.LevelID.String(), Name: v.Level.Name},
			Important:  v.Important,
			Priority:   v.Priority,
			Info:       v.Info,
			ApplicationForm: ApplicationFormDTO{FormSchema: form},
			IsActive:   v.IsActive,
			CreatedAt:  v.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  v.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /admin/vacancies
func AdminCreateVacancy(w http.ResponseWriter, r *http.Request) {
	type VacancyCreateRequest struct {
		Name            string                 `json:"name"`
		LocationID      string                 `json:"locationId"`
		DepartmentID    string                 `json:"departmentId"`
		LevelID         string                 `json:"levelId"`
		Info            string                 `json:"info"`
		ApplicationForm ApplicationFormDTO     `json:"applicationForm"`
		Important       bool                   `json:"important"`
		Priority        int                    `json:"priority"`
		IsActive        bool                   `json:"isActive"`
	}
	var req VacancyCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	locID, err1 := uuid.Parse(req.LocationID)
	depID, err2 := uuid.Parse(req.DepartmentID)
	lvlID, err3 := uuid.Parse(req.LevelID)
	if err1 != nil || err2 != nil || err3 != nil || len(req.Name) < 3 || len(req.Info) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	form, err := json.Marshal(req.ApplicationForm.FormSchema)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid applicationForm"})
		return
	}
	vac := models.Vacancy{
		Name:            req.Name,
		LocationID:      locID,
		DepartmentID:    depID,
		LevelID:         lvlID,
		Info:            req.Info,
		ApplicationForm: datatypes.JSON(form),
		Important:       req.Important,
		Priority:        req.Priority,
		IsActive:        req.IsActive,
	}
	if err := DB.Create(&vac).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedResponse{Success: true, ID: vac.ID.String(), CreatedAt: vac.CreatedAt.Format("2006-01-02T15:04:05Z07:00")})
}

// PATCH /admin/vacancies/{id}
func AdminUpdateVacancy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var v models.Vacancy
	if err := DB.First(&v, "id = ?", id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Resource not found"})
		return
	}
	type VacancyUpdateRequest struct {
		Name            *string                `json:"name"`
		LocationID      *string                `json:"locationId"`
		DepartmentID    *string                `json:"departmentId"`
		LevelID         *string                `json:"levelId"`
		Info            *string                `json:"info"`
		ApplicationForm *ApplicationFormDTO    `json:"applicationForm"`
		Important       *bool                  `json:"important"`
		Priority        *int                   `json:"priority"`
		IsActive        *bool                  `json:"isActive"`
	}
	var req VacancyUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	if req.Name != nil {
		v.Name = *req.Name
	}
	if req.LocationID != nil {
		if locID, err := uuid.Parse(*req.LocationID); err == nil {
			v.LocationID = locID
		}
	}
	if req.DepartmentID != nil {
		if depID, err := uuid.Parse(*req.DepartmentID); err == nil {
			v.DepartmentID = depID
		}
	}
	if req.LevelID != nil {
		if lvlID, err := uuid.Parse(*req.LevelID); err == nil {
			v.LevelID = lvlID
		}
	}
	if req.Info != nil {
		v.Info = *req.Info
	}
	if req.ApplicationForm != nil {
		if form, err := json.Marshal(req.ApplicationForm.FormSchema); err == nil {
			v.ApplicationForm = datatypes.JSON(form)
		}
	}
	if req.Important != nil {
		v.Important = *req.Important
	}
	if req.Priority != nil {
		v.Priority = *req.Priority
	}
	if req.IsActive != nil {
		v.IsActive = *req.IsActive
	}
	if err := DB.Save(&v).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true, "id": v.ID.String()})
}

// DELETE /admin/vacancies/{id}
func AdminDeleteVacancy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := DB.Delete(&models.Vacancy{}, "id = ?", id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /admin/applications
func AdminGetApplications(w http.ResponseWriter, r *http.Request) {
	var apps []models.Application
	query := DB.Preload("Vacancy")
	if status := r.URL.Query().Get("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if vacancyId := r.URL.Query().Get("vacancyId"); vacancyId != "" {
		if vacID, err := uuid.Parse(vacancyId); err == nil {
			query = query.Where("vacancy_id = ?", vacID)
		}
	}
	if err := query.Find(&apps).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	var resp ApplicationsList
	for _, a := range apps {
		var answers map[string]interface{}
		_ = json.Unmarshal(a.Answer, &answers)
		var answeredAt *string
		if a.AnsweredAt != nil {
			ts := a.AnsweredAt.Format("2006-01-02T15:04:05Z07:00")
			answeredAt = &ts
		}
		resp = append(resp, ApplicationItem{
			ID:         a.ID.String(),
			VacancyID:  a.VacancyID.String(),
			VacancyName: a.Vacancy.Name,
			Answers:    answers,
			Status:     a.Status,
			CreatedAt:  a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  a.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
			AnsweredAt: answeredAt,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PATCH /admin/applications/{id}
func AdminUpdateApplication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var app models.Application
	if err := DB.First(&app, "id = ?", id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Resource not found"})
		return
	}
	var req ApplicationStatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request data"})
		return
	}
	app.Status = req.Status
	if err := DB.Save(&app).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"success": true, "id": app.ID.String()})
}

// DELETE /admin/applications/{id}
func AdminDeleteApplication(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := DB.Delete(&models.Application{}, "id = ?", id).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
} 