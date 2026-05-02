package config_test

import (
	"strings"
	"testing"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
)

const (
	envDatabaseURL   = "DATABASE_URL"
	envFirebaseID    = "FIREBASE_PROJECT_ID"
	envGoogleCreds   = "GOOGLE_APPLICATION_CREDENTIALS"
	envAppPort       = "APP_PORT"
	envLogLevel      = "LOG_LEVEL"
	envLogFormat     = "LOG_FORMAT"
	validDatabaseDSN = "postgres://u:p@localhost:5432/db?sslmode=disable"
	validProjectID   = "cooksense-test"
	validCredsPath   = "/secrets/firebase-admin.json"
)

func setAllRequired(t *testing.T) {
	t.Helper()
	t.Setenv(envDatabaseURL, validDatabaseDSN)
	t.Setenv(envFirebaseID, validProjectID)
	t.Setenv(envGoogleCreds, validCredsPath)
}

// TestLoad_AllRequired_ReturnsConfig verifies SPEC-DB-001, SPEC-DB-002, SPEC-DB-003.
func TestLoad_AllRequired_ReturnsConfig(t *testing.T) {
	setAllRequired(t)
	t.Setenv(envAppPort, "9090")
	t.Setenv(envLogLevel, "debug")
	t.Setenv(envLogFormat, "json")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.DatabaseURL != validDatabaseDSN {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, validDatabaseDSN)
	}
	if cfg.FirebaseProjectID != validProjectID {
		t.Errorf("FirebaseProjectID = %q, want %q", cfg.FirebaseProjectID, validProjectID)
	}
	if cfg.GoogleAppCredentials != validCredsPath {
		t.Errorf("GoogleAppCredentials = %q, want %q", cfg.GoogleAppCredentials, validCredsPath)
	}
	if cfg.AppPort != "9090" {
		t.Errorf("AppPort = %q, want %q", cfg.AppPort, "9090")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
	if cfg.LogFormat != "json" {
		t.Errorf("LogFormat = %q, want %q", cfg.LogFormat, "json")
	}
}

// TestLoad_DefaultsApplied verifies optional vars fall back to documented defaults (SPEC-DB-001).
func TestLoad_DefaultsApplied(t *testing.T) {
	setAllRequired(t)
	t.Setenv(envAppPort, "")
	t.Setenv(envLogLevel, "")
	t.Setenv(envLogFormat, "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort default = %q, want %q", cfg.AppPort, "8080")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel default = %q, want %q", cfg.LogLevel, "info")
	}
	if cfg.LogFormat != "text" {
		t.Errorf("LogFormat default = %q, want %q", cfg.LogFormat, "text")
	}
}

// TestLoad_MissingDatabaseURL_ReturnsError verifies SPEC-DB-002.
func TestLoad_MissingDatabaseURL_ReturnsError(t *testing.T) {
	t.Setenv(envDatabaseURL, "")
	t.Setenv(envFirebaseID, validProjectID)
	t.Setenv(envGoogleCreds, validCredsPath)

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load() returned nil error, want error mentioning DATABASE_URL")
	}
	if !strings.Contains(err.Error(), envDatabaseURL) {
		t.Errorf("error = %q, want it to mention %q", err, envDatabaseURL)
	}
}

// TestLoad_MissingFirebaseProjectID_ReturnsError verifies SPEC-DB-002.
func TestLoad_MissingFirebaseProjectID_ReturnsError(t *testing.T) {
	t.Setenv(envDatabaseURL, validDatabaseDSN)
	t.Setenv(envFirebaseID, "")
	t.Setenv(envGoogleCreds, validCredsPath)

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load() returned nil error, want error mentioning FIREBASE_PROJECT_ID")
	}
	if !strings.Contains(err.Error(), envFirebaseID) {
		t.Errorf("error = %q, want it to mention %q", err, envFirebaseID)
	}
}

// TestLoad_MissingGoogleCredentials_ReturnsError verifies SPEC-DB-002.
func TestLoad_MissingGoogleCredentials_ReturnsError(t *testing.T) {
	t.Setenv(envDatabaseURL, validDatabaseDSN)
	t.Setenv(envFirebaseID, validProjectID)
	t.Setenv(envGoogleCreds, "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("Load() returned nil error, want error mentioning GOOGLE_APPLICATION_CREDENTIALS")
	}
	if !strings.Contains(err.Error(), envGoogleCreds) {
		t.Errorf("error = %q, want it to mention %q", err, envGoogleCreds)
	}
}
