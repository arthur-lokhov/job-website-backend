// Package config provides application configuration.
package config

import (
	"os"
	"time"
)

// Config stores all configuration of the application.
// The values are read directly from environment variables.
type Config struct {
	DBUrl      string
	AuthURL    string
	PublicKey  string
	RedisURL   string
	HTTPServer struct {
		Port            string
		Timeout         time.Duration
		ShutdownTimeout time.Duration
	}
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() (config Config, err error) {
	config.DBUrl = os.Getenv("DB_URL")
	config.AuthURL = os.Getenv("AUTH_URL")
	config.PublicKey = os.Getenv("PUBLIC_KEY")
	config.RedisURL = os.Getenv("REDIS_URL")
	config.HTTPServer.Port = os.Getenv("HTTP_SERVER_PORT")

	// Parsing durations from strings
	config.HTTPServer.Timeout, _ = time.ParseDuration(os.Getenv("HTTP_SERVER_TIMEOUT"))
	config.HTTPServer.ShutdownTimeout, _ = time.ParseDuration(os.Getenv("HTTP_SERVER_SHUTDOWN_TIMEOUT"))

	return
}
