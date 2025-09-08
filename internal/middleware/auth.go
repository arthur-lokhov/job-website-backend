// Package middleware provides authentication and authorization middleware for the application.
package middleware

import (
	"context"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"job-website-backend/internal/services"

	"github.com/golang-jwt/jwt/v4"
)

const (
	BearerPrefix             = "Bearer "
	AuthPublicKeyURL         = "https://authservice.cyberzone.dev/v1/auth/publickey?type=pem"
	AuthBlacklistURL         = "https://authservice.cyberzone.dev/v1/auth/expired"
	ClaimExp                 = "exp"
	ClaimNbf                 = "nbf"
	ClaimJti                 = "jti"
	ClaimSub                 = "sub"
	publicKeyRefreshInterval = 5 * time.Minute // Refresh every 5 minutes
	blacklistRefreshInterval = 1 * time.Minute // Refresh every 1 minute
)

var (
	publicKey            interface{} // может быть *rsa.PublicKey или ed25519.PublicKey
	lastPublicKeyRefresh time.Time
	publicKeyMutex       sync.RWMutex
)

type contextKey string

const (
	userContextKey        contextKey = "user"
	permissionsContextKey contextKey = "permissions"
)

func refreshPublicKey() error {
	keyURL := os.Getenv("AUTH_PUBLIC_KEY_URL")
	if keyURL == "" {
		keyURL = os.Getenv("AUTH_URL") + "/v1/auth/publickey?type=pem"
	}
	if keyURL == "" {
		keyURL = "https://authservice.cyberzone.dev/v1/auth/publickey?type=pem"
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, e := client.Get(keyURL)
	if e != nil {
		return e
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	pemBytes, e := io.ReadAll(resp.Body)
	if e != nil {
		return e
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return errors.New("invalid PEM block")
	}
	pub, e := x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		return e
	}

	publicKey = pub
	lastPublicKeyRefresh = time.Now()
	return nil
}

func getPublicKey() (interface{}, error) {
	publicKeyMutex.RLock()
	// Check if key is valid and fresh
	if publicKey != nil && time.Since(lastPublicKeyRefresh) <= publicKeyRefreshInterval {
		defer publicKeyMutex.RUnlock()
		return publicKey, nil
	}
	// Key is not valid or stale, must upgrade to write lock
	publicKeyMutex.RUnlock()

	// Acquire write lock
	publicKeyMutex.Lock()
	defer publicKeyMutex.Unlock()

	// Re-check condition after acquiring write lock, in case another goroutine refreshed it
	// while we were waiting for the lock.
	if publicKey != nil && time.Since(lastPublicKeyRefresh) <= publicKeyRefreshInterval {
		return publicKey, nil
	}

	// Still stale, so we must refresh
	if err := refreshPublicKey(); err != nil {
		return nil, fmt.Errorf("failed to refresh public key: %w", err)
	}

	return publicKey, nil
}

func AuthRequired(authService *services.AuthService, log *slog.Logger, requiredPermissions ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware.AuthRequired"
			log := log.With(slog.String("op", op))

			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
				log.Error("authorization header is missing or invalid")
				unauthorized(w, log)
				return
			}
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

			pubKey, err := getPublicKey()
			if err != nil {
				log.Error("failed to get public key", slog.String("error", err.Error()))
				serverError(w, log)
				return
			}

			claims := jwt.MapClaims{}
			var token *jwt.Token
			if pk, ok := pubKey.(*rsa.PublicKey); ok {
				token, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
					}
					return pk, nil
				})
			} else if pk, ok := pubKey.(ed25519.PublicKey); ok {
				token, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
					}
					return pk, nil
				})
			} else {
				log.Error("unsupported public key type")
				serverError(w, log)
				return
			}

			// Corrected logic:
			if err != nil {
				log.Error("failed to parse or validate token", slog.String("error", err.Error()))
				unauthorized(w, log)
				return
			}

			// Now that we know err is nil, it's safe to check token.Valid
			if !token.Valid {
				log.Error("token is not valid")
				unauthorized(w, log)
				return
			}

			log.Info("token claims", slog.Any("claims", claims))

			// Проверка exp, nbf
			now := time.Now().Unix()
			if !checkTimeClaim(claims, ClaimExp, now, true) || !checkTimeClaim(claims, ClaimNbf, now, false) {
				log.Error("token time claims are invalid")
				unauthorized(w, log)
				return
			}

			// Проверка blacklist через Authservice
			if isBlacklisted(tokenStr) {
				log.Error("token is blacklisted")
				unauthorized(w, log)
				return
			}

			// Получение и кэширование прав пользователя
			uid, _ := claims["sub"].(string)
			var perms []string
			if uid != "" {
				perms, _ = authService.GetUserPermissions(r.Context(), uid)
			}

			log.Info("user permissions", slog.Any("permissions", perms))

			// Check if user has all required permissions
			if len(requiredPermissions) > 0 && !checkPermissions(perms, requiredPermissions) {
				log.Warn("user does not have required permissions", slog.String("user_id", uid), slog.Any("required_perms", requiredPermissions), slog.Any("user_perms", perms))
				forbidden(w, log)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, claims)
			ctx = context.WithValue(ctx, permissionsContextKey, perms)
			next(w, r.WithContext(ctx))
		}
	}
}

func checkTimeClaim(claims jwt.MapClaims, key string, now int64, mustBeFuture bool) bool {
	v, ok := claims[key]
	slog.Default().Info("Checking time claim", slog.String("key", key), slog.Any("value", v), slog.Bool("ok", ok))
	if !ok {
		return false
	}
	var t int64
	switch val := v.(type) {
	case float64:
		t = int64(val)
	case int64:
		t = val
	default:
		return false
	}
	if mustBeFuture {
		return now < t
	}
	return now >= t
}

var (
	blacklistedJTIs      map[string]struct{}
	lastBlacklistRefresh time.Time
	blacklistMutex       sync.RWMutex
)

func refreshBlacklist() error {
	url := os.Getenv("AUTH_URL") + "/v1/auth/expired"
	if url == "/v1/auth/expired" {
		url = "https://authservice.cyberzone.dev/v1/auth/expired"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create blacklist request: %w", err)
	}

	apiToken := os.Getenv("AUTH_API_TOKEN")
	if apiToken != "" {
		req.Header.Set("Api-Token", apiToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch blacklist: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch blacklist: status code %d", resp.StatusCode)
	}

	var jtis []string
	if err := json.NewDecoder(resp.Body).Decode(&jtis); err != nil {
		return fmt.Errorf("failed to decode blacklist: %w", err)
	}

	newBlacklist := make(map[string]struct{})
	for _, jti := range jtis {
		newBlacklist[jti] = struct{}{}
	}

	blacklistedJTIs = newBlacklist
	lastBlacklistRefresh = time.Now()
	return nil
}

func isBlacklisted(token string) bool {
	claims := jwt.MapClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(token, claims)
	if err != nil {
		return false // Cannot parse token, treat as not blacklisted
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return false // No JTI claim, treat as not blacklisted
	}

	blacklistMutex.RLock()
	if blacklistedJTIs != nil && time.Since(lastBlacklistRefresh) <= blacklistRefreshInterval {
		_, found := blacklistedJTIs[jti]
		blacklistMutex.RUnlock()
		return found
	}
	blacklistMutex.RUnlock()

	blacklistMutex.Lock()
	defer blacklistMutex.Unlock()

	if blacklistedJTIs != nil && time.Since(lastBlacklistRefresh) <= blacklistRefreshInterval {
		_, found := blacklistedJTIs[jti]
		return found
	}

	if err := refreshBlacklist(); err != nil {
		fmt.Printf("failed to refresh blacklist: %v\n", err)
		return true // Fail closed
	}

	_, found := blacklistedJTIs[jti]
	return found
}

func unauthorized(w http.ResponseWriter, log *slog.Logger) {
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"}); err != nil {
		log.Error("failed to write unauthorized response", slog.String("error", err.Error()))
	}
}

func forbidden(w http.ResponseWriter, log *slog.Logger) {
	w.WriteHeader(http.StatusForbidden)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"}); err != nil {
		log.Error("failed to write forbidden response", slog.String("error", err.Error()))
	}
}

func serverError(w http.ResponseWriter, log *slog.Logger) {
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"}); err != nil {
		log.Error("failed to write server error response", slog.String("error", err.Error()))
	}
}

// checkPermissions checks if the user's permissions contain all the required permissions.
func checkPermissions(userPermissions []string, requiredPermissions []string) bool {
	for _, reqPerm := range requiredPermissions {
		found := false
		for _, userPerm := range userPermissions {
			if userPerm == reqPerm {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
