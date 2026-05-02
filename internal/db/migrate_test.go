//go:build integration

package db_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlefebvre1/cooksense-backend/internal/db"
)

// TestUp_IsIdempotent verifies SPEC-DB-008, SPEC-DB-009.
func TestUp_IsIdempotent(t *testing.T) {
	dsn := testDSN(t)
	dir := migrationsPath(t)

	if err := db.Up(t.Context(), dsn, dir); err != nil {
		t.Fatalf("first Up: %v", err)
	}
	if err := db.Up(t.Context(), dsn, dir); err != nil {
		t.Fatalf("second Up (should be idempotent): %v", err)
	}
}

// TestDown_OneStep_RollsBack verifies SPEC-DB-010.
func TestDown_OneStep_RollsBack(t *testing.T) {
	dsn := testDSN(t)
	dir := migrationsPath(t)

	if err := db.Up(t.Context(), dsn, dir); err != nil {
		t.Fatalf("Up: %v", err)
	}
	if err := db.Down(t.Context(), dsn, dir, 1); err != nil {
		t.Fatalf("Down(1): %v", err)
	}
	// Restore schema for subsequent integration tests in the same run.
	if err := db.Up(t.Context(), dsn, dir); err != nil {
		t.Fatalf("Up (restore): %v", err)
	}
}

// TestDown_NonPositiveN_ReturnsError verifies SPEC-DB-010 input validation.
func TestDown_NonPositiveN_ReturnsError(t *testing.T) {
	dsn := testDSN(t)
	dir := migrationsPath(t)
	if err := db.Down(t.Context(), dsn, dir, 0); err == nil {
		t.Error("Down(0) returned nil, want error")
	}
	if err := db.Down(t.Context(), dsn, dir, -3); err == nil {
		t.Error("Down(-3) returned nil, want error")
	}
}

// TestUp_BadSQL_WrapsMigrationError verifies SPEC-DB-011: when the migration
// runner reports a failure, db.Up wraps it in *MigrationError carrying the
// version of the failing migration. Uses an isolated tempdir so the real
// schema is not touched.
func TestUp_BadSQL_WrapsMigrationError(t *testing.T) {
	dsn := testDSN(t)
	// Run in an isolated schema so we don't disturb the canonical public-schema
	// state used by the other integration tests in this package.
	schema := newIsolatedSchema(t, dsn)
	scopedDSN := withSearchPath(t, dsn, schema)

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "0001_bad.up.sql"), []byte("THIS IS NOT SQL;"), 0o600); err != nil {
		t.Fatalf("write up sql: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "0001_bad.down.sql"), []byte("-- noop"), 0o600); err != nil {
		t.Fatalf("write down sql: %v", err)
	}

	err := db.Up(t.Context(), scopedDSN, dir)
	if err == nil {
		t.Fatal("Up against bad SQL returned nil, want MigrationError")
	}
	var me *db.MigrationError
	if !errors.As(err, &me) {
		t.Fatalf("errors.As to *MigrationError = false, got %T: %v", err, err)
	}
}

// TestMigrationError_WrapsVersion verifies SPEC-DB-011.
func TestMigrationError_WrapsVersion(t *testing.T) {
	inner := errors.New("syntax error at or near \"FOO\"")
	err := &db.MigrationError{Version: 7, Err: inner}

	if !strings.Contains(err.Error(), "7") {
		t.Errorf("Error() = %q, want it to contain version 7", err.Error())
	}
	if !errors.Is(err, inner) {
		t.Errorf("errors.Is(err, inner) = false, want true (Unwrap broken)")
	}
	var target *db.MigrationError
	if !errors.As(err, &target) {
		t.Fatal("errors.As to *MigrationError failed")
	}
	if target.Version != 7 {
		t.Errorf("target.Version = %d, want 7", target.Version)
	}
}
