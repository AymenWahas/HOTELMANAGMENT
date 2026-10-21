package config

import (
	"fmt"
	"os"
)

type Config struct{ Port, DatabaseURL, APIToken, Timezone string }

func Load() (Config, error) {
	cfg := Config{Port: os.Getenv("APP_PORT"), DatabaseURL: os.Getenv("DATABASE_URL"), APIToken: os.Getenv("API_TOKEN"), Timezone: os.Getenv("HOTEL_TIMEZONE")}
	if cfg.Port == "" {
		cfg.Port = "8089"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.APIToken == "" {
		return Config{}, fmt.Errorf("API_TOKEN is required")
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Riyadh"
	}
	return cfg, nil
}
