package config

import (
	"fmt"

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

	return cfg, nil
}
