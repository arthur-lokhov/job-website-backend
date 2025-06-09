package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"git.cyberzone.dev/project-vacancy-website/job-website-backend/internal/services"
)

type contextKey string

const (
	UserIDKey      contextKey = "user_id"
	UserEmailKey   contextKey = "user_email"
	UserRolesKey   contextKey = "user_roles"
	UserTokenKey   contextKey = "user_token"
	UserClaimsKey  contextKey = "user_claims"
)

func Auth(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := authService.ExtractTokenFromHeader(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := authService.ValidateToken(r.Context(), tokenString)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Extract user information from claims
			userID, _ := uuid.Parse(claims["sub"].(string))
			email := claims["email"].(string)
			roles := claims["roles"].([]interface{})

			// Create a new context with user information
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserEmailKey, email)
			ctx = context.WithValue(ctx, UserRolesKey, roles)
			ctx = context.WithValue(ctx, UserTokenKey, token)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value(UserTokenKey).(*jwt.Token)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			authService := r.Context().Value("auth_service").(*services.AuthService)
			if err := authService.CheckPermission(r.Context(), token, permission); err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions to get user information from context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

func GetUserEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

func GetUserRoles(ctx context.Context) ([]interface{}, bool) {
	roles, ok := ctx.Value(UserRolesKey).([]interface{})
	return roles, ok
}

func GetUserToken(ctx context.Context) (*jwt.Token, bool) {
	token, ok := ctx.Value(UserTokenKey).(*jwt.Token)
	return token, ok
}

func GetUserClaims(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(jwt.MapClaims)
	return claims, ok
} 