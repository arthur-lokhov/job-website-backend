// Package services contains the business logic of the application.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"job-website-backend/internal/dto"

	"job-website-backend/internal/repository"

	"github.com/go-redis/redis/v8"
)

const (
	RedisKeyDepartments = "departments"
)

type DepartmentService struct {
	repo        *repository.DepartmentRepository
	authURL     string
	redisClient *redis.Client
	log         *slog.Logger
}

func NewDepartmentService(repo *repository.DepartmentRepository, authURL string, redisClient *redis.Client, log *slog.Logger) *DepartmentService {
	return &DepartmentService{repo: repo, authURL: authURL, redisClient: redisClient, log: log}
}

func (s *DepartmentService) GetDepartments(ctx context.Context) ([]dto.DepartmentDTO, error) {
	const op = "services.DepartmentService.GetDepartments"
	log := s.log.With(slog.String("op", op))

	log.Info("getting departments")
	// Try to get from cache first
	val, err := s.redisClient.Get(ctx, RedisKeyDepartments).Result()
	if err == redis.Nil {
		log.Info("cache miss, fetching from auth service")
		depts, err := s.FetchAndCacheDepartments()
		if err != nil {
			log.Error("failed to fetch and cache departments", slog.String("error", err.Error()))
			return nil, err
		}
		return depts, nil
	} else if err != nil {
		log.Error("failed to get departments from redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get departments from redis: %w", err)
	}

	log.Info("cache hit, returning departments")
	var depts []dto.DepartmentDTO
	if err := json.Unmarshal([]byte(val), &depts); err != nil {
		log.Error("failed to unmarshal departments from redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to unmarshal departments from redis: %w", err)
	}

	return depts, nil
}

func (s *DepartmentService) FetchAndCacheDepartments() ([]dto.DepartmentDTO, error) {
	const op = "services.DepartmentService.FetchAndCacheDepartments"
	log := s.log.With(slog.String("op", op))

	log.Info("fetching and caching departments")
	req, err := http.NewRequest("GET", s.authURL+"/v1/departments", nil)
	if err != nil {
		log.Error("failed to create request", slog.String("error", err.Error()))
		return nil, err
	}

	apiToken := os.Getenv("AUTH_API_TOKEN") // Get API token
	if apiToken != "" {
		req.Header.Set("Api-Token", apiToken) // Set Api-Token header
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("failed to do request", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Error("failed to fetch departments from authservice", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to fetch departments from authservice: status code %d", resp.StatusCode)
	}

	var departments []dto.DepartmentDTO
	if err := json.NewDecoder(resp.Body).Decode(&departments); err != nil {
		log.Error("failed to decode response body", slog.String("error", err.Error()))
		return nil, err
	}

	if err := s.repo.UpsertDepartments(context.Background(), departments); err != nil {
		log.Error("failed to upsert departments", slog.String("error", err.Error()))
		// We can choose to continue even if this fails, as the data will be in cache
	}

	log.Info("caching departments")
	// Marshal departments to JSON before storing in Redis
	jsonData, err := json.Marshal(departments)
	if err != nil {
		log.Error("failed to marshal departments to JSON", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to marshal departments to JSON: %w", err)
	}

	// Store in Redis with a 1-minute expiration
	if err := s.redisClient.Set(context.Background(), RedisKeyDepartments, jsonData, 1*time.Minute).Err(); err != nil {
		log.Error("failed to set departments in redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to set departments in redis: %w", err)
	}

	log.Info("successfully fetched and cached departments")
	return departments, nil
}
