// Package config provides runtime configuration loading and validation for cooksense-server.
package config

import (
	"fmt"
	"os"
)

// Config holds the validated runtime configuration for cooksense-server.
//
// SPEC-DB-001, SPEC-DB-003
type Config struct {
	AppPort              string
	LogLevel             string
	LogFormat            string
	DatabaseURL          string
	FirebaseProjectID    string
	GoogleAppCredentials string
}

// Load reads configuration from environment variables and validates required fields.
// It returns an error naming each missing required variable.
//
// SPEC-DB-001, SPEC-DB-002, SPEC-DB-003
func Load() (Config, error) {
	cfg := Config{
		AppPort:              getEnvDefault("APP_PORT", "8080"),
		LogLevel:             getEnvDefault("LOG_LEVEL", "info"),
		LogFormat:            getEnvDefault("LOG_FORMAT", "text"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		FirebaseProjectID:    os.Getenv("FIREBASE_PROJECT_ID"),
		GoogleAppCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.FirebaseProjectID == "" {
		missing = append(missing, "FIREBASE_PROJECT_ID")
	}
	if cfg.GoogleAppCredentials == "" {
		missing = append(missing, "GOOGLE_APPLICATION_CREDENTIALS")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("config: missing required environment variables: %v", missing)
	}
	return cfg, nil
}

func getEnvDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
