package services

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	publicKeyURL string
	publicKey    *rsa.PublicKey
	client       *http.Client
}

func NewAuthService(publicKeyURL string) *AuthService {
	return &AuthService{
		publicKeyURL: publicKeyURL,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	if s.publicKey == nil {
		if err := s.fetchPublicKey(ctx); err != nil {
			return nil, fmt.Errorf("failed to fetch public key: %w", err)
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, ErrUnauthorized
	}

	// Check if token is blacklisted
	if err := s.checkBlacklist(ctx, tokenString); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) CheckPermission(ctx context.Context, token *jwt.Token, permission string) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrUnauthorized
	}

	permissions, ok := claims["permissions"].([]interface{})
	if !ok {
		return ErrForbidden
	}

	for _, p := range permissions {
		if p == permission {
			return nil
		}
	}

	return ErrForbidden
}

func (s *AuthService) fetchPublicKey(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", s.publicKeyURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch public key: status %d", resp.StatusCode)
	}

	var result struct {
		PublicKey string `json:"public_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(result.PublicKey))
	if err != nil {
		return err
	}

	s.publicKey = key
	return nil
}

func (s *AuthService) checkBlacklist(ctx context.Context, tokenString string) error {
	// Extract token ID from the token
	token, err := jwt.Parse(tokenString, nil)
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrUnauthorized
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return ErrUnauthorized
	}

	// Check if token is blacklisted
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/blacklist/%s", strings.TrimSuffix(s.publicKeyURL, "/public-key"), jti), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ErrUnauthorized
	}

	return nil
}

func (s *AuthService) ExtractTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", ErrUnauthorized
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrUnauthorized
	}

	return parts[1], nil
} 