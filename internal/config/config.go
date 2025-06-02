package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config structure for storing all settings
type Config struct {
	App struct {
		Name  string `mapstructure:"APP_NAME"`   // The name of the application.
		IsDev bool   `mapstructure:"APP_IS_DEV"` // The current environment(IsProd or IsDev).
		Port  string `mapstructure:"APP_PORT"`   // The port on which the application will run.
	}

	Database struct {
		URL      string `mapstructure:"DB_URL"`      // The database connection URL (used for convenience or overriding other parameters).
		User     string `mapstructure:"DB_USER"`     // The username for the database connection.
		Password string `mapstructure:"DB_PASSWORD"` // The password for the database connection.
		Name     string `mapstructure:"DB_NAME"`     // The name of the database to connect to.
		Port     int    `mapstructure:"DB_PORT"`     // The port number of the database server.
		Host     string `mapstructure:"DB_HOST"`     // The hostname or IP address of the database server.
		SSLMode  string `mapstructure:"DB_SSLMODE"`  // The SSL mode for the database connection (e.g., disable, require, verify-full).
	}

	Logger struct {
		Level string `mapstructure:"LOG_LEVEL"` // The logging level (e.g., debug, info, warn, error, fatal).
	}
}

func LoadConfig(path string, logger *zap.Logger) (*Config, error) {
	viper.SetConfigFile(fmt.Sprintf("%s/.env", path))
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		logger.Info("No .env file found, using environment variables only")
	}

	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("APP_NAME", "myapp")
	viper.SetDefault("APP_IS_DEV", true)
	viper.SetDefault("APP_PORT", "8080")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_NAME", "mydb")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_SSLMODE", "disable")

	viper.SetDefault("LOG_LEVEL", "debug")
}

func (c *Config) IsProduction() bool {
	return !c.App.IsDev
}
