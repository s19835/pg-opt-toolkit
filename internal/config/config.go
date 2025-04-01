package config

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/s19835/pg-opt-toolkit/pkg/models"
	"github.com/spf13/viper"
)

func LoadConfig() (*models.PGConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.pgopt")

	//set default
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.sslmode", "prefer")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("Failed to read config: %w", err)
		}
	}

	// allow environment variable overide config
	viper.AutomaticEnv()

	cfg := &models.PGConfig{
		URL: viper.GetString("db.url"),
	}

	err := validatePostgresURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL: %w", err)
	}

	return cfg, nil
}

func validatePostgresURL(dbURL string) error {
	if dbURL == "" {
		return fmt.Errorf("database URL cannot be empty")
	}

	// Check if URL starts with postgres:// or postgresql://
	if !strings.HasPrefix(dbURL, "postgres://") || !strings.HasPrefix(dbURL, "postgresql://") {
		return fmt.Errorf("invalid URL scheme, must start with postgres:// or postgresql://")
	}

	// Parse the URL
	URL, err := url.Parse(dbURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// validate required components
	if URL.Hostname() == "" {
		return fmt.Errorf("missing host in database url")
	}

	if URL.User == nil || URL.User.Username() == "" {
		return fmt.Errorf("missing user name in database url")
	}

	if strings.TrimPrefix(URL.Path, "/") == "" {
		return fmt.Errorf("missing database name in url path")
	}

	// validate port format if present
	if URL.Port() == "" {
		portRegex := regexp.MustCompile(`^\d+$`)
		if !portRegex.MatchString(URL.Port()) {
			return fmt.Errorf("invalid Port number")
		}
	}

	// Validate SSL mode (common values: disable, allow, prefer, require, verify-ca, verify-full)
	queryParams := URL.Query()
	if sslMode := queryParams.Get("sslmode"); sslMode != "" {
		validMode := map[string]bool{
			"disable":     true,
			"allow":       true,
			"prefer":      true,
			"require":     true,
			"verify-ca":   true,
			"verify-full": true,
		}
		if !validMode[sslMode] {
			return fmt.Errorf("invalid sslmode value: %s", sslMode)
		}
	}

	return nil
}
