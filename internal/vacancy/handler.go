package vacancy

import (
	"encoding/json"
	"log"
	"net/http"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/db"
	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/models"

	"github.com/google/uuid"
)

func CreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancy
	if err := json.NewDecoder(r.Body).Decode(&vacancy); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	vacancy.ID = uuid.New()

	if err := db.DB.Create(&vacancy).Error; err != nil {
		log.Printf("error creating vacancy: %v", err)
		http.Error(w, "Failed to create vacancy", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vacancy)
}

func GetAllVacancies(w http.ResponseWriter, r *http.Request) {
	var vacancies []models.Vacancy
	if err := db.DB.Preload("Department").Preload("Level").Preload("Location").Find(&vacancies).Error; err != nil {
		http.Error(w, "Failed to get vacancies", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vacancies)
}
