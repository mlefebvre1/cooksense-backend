//go:build integration

package db_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
	"github.com/mlefebvre1/cooksense-backend/internal/db"
)

// testDB opens a real pool against DATABASE_URL, skips if DATABASE_URL is unset,
// and registers t.Cleanup to close the pool.
//
// SPEC-DB-09 §8.3
func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	pool, err := db.Open(t.Context(), config.Config{DatabaseURL: dsn})
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func testDSN(t *testing.T) string {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	return dsn
}

func migrationsPath(t *testing.T) string {
	t.Helper()
	// Tests run from internal/db; migrations live two levels up.
	return "../../migrations"
}

// newIsolatedSchema creates a fresh, throw-away Postgres schema and registers
// a t.Cleanup that drops it. Used by tests that need to run migrations
// without touching the canonical public-schema state.
func newIsolatedSchema(t *testing.T, dsn string) string {
	t.Helper()
	name := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), rand.IntN(1_000_000))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("newIsolatedSchema: connect: %v", err)
	}
	defer conn.Close(ctx)
	if _, err := conn.Exec(ctx, fmt.Sprintf("CREATE SCHEMA %q", name)); err != nil {
		t.Fatalf("newIsolatedSchema: create: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pgx.Connect(ctx, dsn)
		if err != nil {
			t.Logf("newIsolatedSchema cleanup: connect: %v", err)
			return
		}
		defer conn.Close(ctx)
		if _, err := conn.Exec(ctx, fmt.Sprintf("DROP SCHEMA IF EXISTS %q CASCADE", name)); err != nil {
			t.Logf("newIsolatedSchema cleanup: drop: %v", err)
		}
	})
	return name
}

// withSearchPath returns dsn with search_path appended (or replaced) so all
// statements run by this connection target the given schema instead of public.
func withSearchPath(t *testing.T, dsn, schema string) string {
	t.Helper()
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("withSearchPath: parse dsn: %v", err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	return u.String()
}
