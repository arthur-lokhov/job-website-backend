// Package services contains the business logic of the application.
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type AuthService struct {
	redisClient *redis.Client
	log         *slog.Logger
}

func NewAuthService(redisClient *redis.Client, log *slog.Logger) *AuthService {
	return &AuthService{redisClient: redisClient, log: log}
}

func (s *AuthService) GetUserPermissions(ctx context.Context, uid string) ([]string, error) {
	const op = "services.AuthService.GetUserPermissions"
	log := s.log.With(slog.String("op", op), slog.String("uid", uid))

	log.Info("getting user permissions")

	val, err := s.redisClient.Get(ctx, uid).Result()
	if err == redis.Nil {
		log.Info("cache miss, fetching from auth service")
		// Cache miss, fetch from auth service
		authURL := os.Getenv("AUTH_URL")
		if authURL == "" {
			authURL = "https://authservice.cyberzone.dev"
		}

		req, err := http.NewRequest("GET", authURL+"/v1/admin/permissions?uuid="+uid, nil)
		if err != nil {
			log.Error("failed to create request", slog.String("error", err.Error()))
			return nil, err
		}
		apiToken := os.Getenv("AUTH_API_TOKEN")
		if apiToken != "" {
			req.Header.Set("Api-Token", apiToken)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error("failed to do request", slog.String("error", err.Error()))
			return nil, err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Error("failed to close response body", slog.String("error", err.Error()))
			}
		}()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Error("failed to fetch permissions: non-OK status",
				slog.Int("status_code", resp.StatusCode),
				slog.String("response_body", string(bodyBytes)))
			return nil, fmt.Errorf("failed to fetch permissions: status %d", resp.StatusCode)
		}

		var permissions []string
		if err := json.NewDecoder(resp.Body).Decode(&permissions); err != nil {
			log.Error("failed to decode response body", slog.String("error", err.Error()))
			return nil, err
		}

		log.Info("caching permissions")
		// Cache the permissions
		jsonData, err := json.Marshal(permissions)
		if err != nil {
			log.Error("failed to marshal permissions to JSON", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to marshal permissions to JSON: %w", err)
		}
		if err := s.redisClient.Set(ctx, uid, jsonData, 24*time.Hour).Err(); err != nil {
			log.Error("failed to set permissions in redis", slog.String("error", err.Error()))
			return nil, fmt.Errorf("failed to set permissions in redis: %w", err)
		}

		log.Info("successfully fetched and cached permissions")
		return permissions, nil
	} else if err != nil {
		log.Error("failed to get permissions from redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get permissions from redis: %w", err)
	}

	log.Info("cache hit, returning permissions")
	// Cache hit, unmarshal and return
	var permissions []string
	if err := json.Unmarshal([]byte(val), &permissions); err != nil {
		log.Error("failed to unmarshal permissions from redis", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to unmarshal permissions from redis: %w", err)
	}

	return permissions, nil
}
