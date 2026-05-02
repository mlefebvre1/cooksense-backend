package db

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
)

// TestPgxMigrateURL_ConvertsSchemes verifies the pgx5:// scheme conversion
// used by db.Up and db.Down. SPEC-DB-008, SPEC-DB-010.
func TestPgxMigrateURL_ConvertsSchemes(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"postgres scheme", "postgres://u:p@host:5432/db?sslmode=disable", "pgx5://u:p@host:5432/db?sslmode=disable"},
		{"postgresql scheme", "postgresql://u:p@host:5432/db?sslmode=disable", "pgx5://u:p@host:5432/db?sslmode=disable"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pgxMigrateURL(tc.in)
			if err != nil {
				t.Fatalf("pgxMigrateURL(%q) error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("pgxMigrateURL(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestPgxMigrateURL_RejectsUnknownScheme verifies SPEC-DB-008 input validation.
func TestPgxMigrateURL_RejectsUnknownScheme(t *testing.T) {
	_, err := pgxMigrateURL("mysql://u:p@host/db")
	if err == nil {
		t.Fatal("expected error for unsupported scheme, got nil")
	}
	if !strings.Contains(err.Error(), "postgres://") {
		t.Errorf("error = %q, want it to mention postgres://", err)
	}
}

// TestUp_InvalidDSN_ReturnsError verifies SPEC-DB-008 input validation: an
// unparseable DSN must surface as an error before any I/O is attempted.
func TestUp_InvalidDSN_ReturnsError(t *testing.T) {
	if err := Up(context.Background(), "mysql://nope", t.TempDir()); err == nil {
		t.Fatal("Up with invalid scheme returned nil, want error")
	}
}

// TestDown_InvalidDSN_ReturnsError verifies SPEC-DB-010 input validation.
func TestDown_InvalidDSN_ReturnsError(t *testing.T) {
	if err := Down(context.Background(), "mysql://nope", t.TempDir(), 1); err == nil {
		t.Fatal("Down with invalid scheme returned nil, want error")
	}
}

// TestDown_NonPositiveN_ReturnsError verifies SPEC-DB-010: n <= 0 is rejected
// without touching the database.
func TestDown_NonPositiveN_ReturnsError(t *testing.T) {
	for _, n := range []int{0, -1, -42} {
		err := Down(context.Background(), "postgres://u:p@host/db", "migrations", n)
		if err == nil {
			t.Errorf("Down(n=%d) returned nil, want error", n)
			continue
		}
		if !strings.Contains(err.Error(), "positive") {
			t.Errorf("Down(n=%d) error = %q, want mention of 'positive'", n, err)
		}
	}
}

// TestOpen_InvalidDSN_ReturnsError verifies SPEC-DB-004 input validation: a
// DSN that pgxpool.ParseConfig cannot parse must surface as an error.
func TestOpen_InvalidDSN_ReturnsError(t *testing.T) {
	cfg := config.Config{DatabaseURL: "::not a dsn::"}
	pool, err := Open(context.Background(), cfg)
	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("Open with invalid DSN returned nil error, want error")
	}
}

// TestOpen_PingFailure_ClosesPoolAndErrors verifies SPEC-DB-006: a successful
// parse but failing ping must close the pool and propagate the error.
func TestOpen_PingFailure_ClosesPoolAndErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	cfg := config.Config{
		DatabaseURL: "postgres://u:p@127.0.0.1:1/db?sslmode=disable&connect_timeout=1",
	}
	pool, err := Open(ctx, cfg, func(c *pgxpool.Config) {
		c.MinConns = 0
	})
	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("Open against unreachable host returned nil error, want ping failure")
	}
}

// TestOpen_OverrideApplied verifies SPEC-DB-007: the variadic override is
// applied to the pool config before pgxpool.NewWithConfig.
func TestOpen_OverrideApplied(t *testing.T) {
	called := false
	cfg := config.Config{DatabaseURL: "postgres://u:p@127.0.0.1:1/db?sslmode=disable&connect_timeout=1"}
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	pool, _ := Open(ctx, cfg, func(c *pgxpool.Config) {
		called = true
		c.MaxConns = 5
	})
	if pool != nil {
		pool.Close()
	}
	if !called {
		t.Error("override callback was not invoked")
	}
}

// TestMigrationError_Format verifies SPEC-DB-011: Error() reports the version
// and Unwrap() exposes the wrapped error so errors.Is/As work.
func TestMigrationError_Format(t *testing.T) {
	inner := errors.New("boom")
	e := &MigrationError{Version: 42, Err: inner}
	if !strings.Contains(e.Error(), "42") {
		t.Errorf("Error() = %q, want it to contain version 42", e.Error())
	}
	if !errors.Is(e, inner) {
		t.Error("errors.Is(e, inner) = false, want true")
	}
}
