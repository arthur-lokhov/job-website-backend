package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
)

var (
	publicKey     interface{} // может быть *rsa.PublicKey или ed25519.PublicKey
	publicKeyOnce sync.Once
)

func getPublicKey() (interface{}, error) {
	var err error
	publicKeyOnce.Do(func() {
		keyUrl := os.Getenv("AUTH_PUBLIC_KEY_URL")
		if keyUrl == "" {
			keyUrl = os.Getenv("AUTH_URL") + "/v1/auth/publickey"
		}
		if keyUrl == "" {
			keyUrl = "https://authservice.cyberzone.dev/v1/auth/publickey"
		}
		resp, e := http.Get(keyUrl)
		if e != nil {
			err = e
			return
		}
		defer resp.Body.Close()
		pemBytes, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			err = e
			return
		}
		block, _ := pem.Decode(pemBytes)
		if block == nil {
			err = errors.New("invalid PEM block")
			return
		}
		pub, e := x509.ParsePKIXPublicKey(block.Bytes)
		if e != nil {
			err = e
			return
		}
		publicKey = pub
	})
	return publicKey, err
}

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			unauthorized(w)
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		pubKey, err := getPublicKey()
		if err != nil {
			serverError(w)
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
			serverError(w)
			return
		}
		if err != nil || !token.Valid {
			unauthorized(w)
			return
		}
		// Проверка exp, nbf
		now := time.Now().Unix()
		if !checkTimeClaim(claims, "exp", now, true) || !checkTimeClaim(claims, "nbf", now, false) {
			unauthorized(w)
			return
		}
		// Проверка blacklist через Authservice
		if isBlacklisted(tokenStr) {
			unauthorized(w)
			return
		}
		// Получение и кэширование прав пользователя
		uid, _ := claims["uid"].(string)
		var perms []string
		if uid != "" {
			perms, _ = GetUserPermissions(uid)
		}
		ctx := context.WithValue(r.Context(), "user", claims)
		ctx = context.WithValue(ctx, "permissions", perms)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkTimeClaim(claims jwt.MapClaims, key string, now int64, mustBeFuture bool) bool {
	v, ok := claims[key]
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

func isBlacklisted(token string) bool {
	url := os.Getenv("AUTH_URL") + "/v1/auth/expired?token=" + token
	if url == "/v1/auth/expired?token=" {
		url = "https://authservice.cyberzone.dev/v1/auth/expired?token=" + token
	}
	resp, err := http.Get(url)
	if err != nil {
		return false // fail open
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var res struct {
			Expired bool `json:"expired"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		return res.Expired
	}
	return false
}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "Authentication required"})
}

func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
}

