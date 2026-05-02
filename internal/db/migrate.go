package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrationError wraps a failure from the migration runner together with the
// version number of the migration that failed, for easier debugging.
//
// SPEC-DB-011
type MigrationError struct {
	Version uint
	Err     error
}

// Error implements the error interface.
func (e *MigrationError) Error() string {
	return fmt.Sprintf("db: migration version %d failed: %v", e.Version, e.Err)
}

// Unwrap returns the underlying error so callers can use errors.Is / errors.As.
func (e *MigrationError) Unwrap() error {
	return e.Err
}

// Up applies all pending up migrations from migrationsDir against the database
// identified by dsn. It treats migrate.ErrNoChange as success, so calling Up
// when all migrations are already applied returns nil.
//
// SPEC-DB-008, SPEC-DB-009, SPEC-DB-011
func Up(_ context.Context, dsn, migrationsDir string) error {
	m, err := newMigrator(dsn, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return wrapMigrationError(m, err)
	}
	return nil
}

// Down rolls back exactly n migration steps from migrationsDir against dsn.
// Passing n <= 0 returns an error.
//
// SPEC-DB-010, SPEC-DB-011
func Down(_ context.Context, dsn, migrationsDir string, n int) error {
	if n <= 0 {
		return errors.New("db: n must be positive")
	}
	m, err := newMigrator(dsn, migrationsDir)
	if err != nil {
		return err
	}
	defer closeMigrator(m)

	if err := m.Steps(-n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return wrapMigrationError(m, err)
	}
	return nil
}

func newMigrator(dsn, migrationsDir string) (*migrate.Migrate, error) {
	sourceURL := "file://" + migrationsDir
	databaseURL, err := pgxMigrateURL(dsn)
	if err != nil {
		return nil, err
	}
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: open migrator: %w", err)
	}
	return m, nil
}

// pgxMigrateURL converts a postgres:// DSN into the pgx5:// scheme expected by
// the golang-migrate pgx/v5 driver.
func pgxMigrateURL(dsn string) (string, error) {
	const pgxPrefix = "pgx5://"
	switch {
	case strings.HasPrefix(dsn, "postgres://"):
		return pgxPrefix + strings.TrimPrefix(dsn, "postgres://"), nil
	case strings.HasPrefix(dsn, "postgresql://"):
		return pgxPrefix + strings.TrimPrefix(dsn, "postgresql://"), nil
	default:
		return "", errors.New("db: DATABASE_URL must start with postgres:// or postgresql://")
	}
}

func wrapMigrationError(m *migrate.Migrate, err error) error {
	version, _, vErr := m.Version()
	if vErr != nil {
		version = 0
	}
	return &MigrationError{Version: version, Err: err}
}

func closeMigrator(m *migrate.Migrate) {
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		slog.Warn("db: migrator source close error", "err", srcErr)
	}
	if dbErr != nil {
		slog.Warn("db: migrator database close error", "err", dbErr)
	}
}
