//go:build integration

package users_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlefebvre1/cooksense-backend/internal/auth"
	"github.com/mlefebvre1/cooksense-backend/internal/config"
	"github.com/mlefebvre1/cooksense-backend/internal/db"
	"github.com/mlefebvre1/cooksense-backend/internal/users"
)

// testPool opens a pgx pool against DATABASE_URL (skips if unset) and applies
// the canonical migrations so the users table exists. Returns a closed pool
// at end-of-test via t.Cleanup.
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set")
	}
	if err := db.Up(t.Context(), dsn, "../../migrations"); err != nil {
		t.Fatalf("db.Up: %v", err)
	}
	pool, err := db.Open(t.Context(), config.Config{DatabaseURL: dsn})
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// TestRepo_Touch_InsertsThenUpdates verifies SPEC-AUTH-018, SPEC-AUTH-019:
// the canonical UPSERT inserts a fresh row and updates last_seen_at on conflict.
func TestRepo_Touch_InsertsThenUpdates(t *testing.T) {
	pool := testPool(t)
	repo := users.NewRepo(pool)
	uid := "test-uid-" + time.Now().Format("150405.000000000")

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE firebase_uid = $1", uid)
	})

	u := auth.User{UID: uid, Email: "x@example.com", DisplayName: "X"}
	if err := repo.Touch(t.Context(), u); err != nil {
		t.Fatalf("Touch (insert): %v", err)
	}
	var firstSeen time.Time
	if err := pool.QueryRow(t.Context(), "SELECT last_seen_at FROM users WHERE firebase_uid = $1", uid).Scan(&firstSeen); err != nil {
		t.Fatalf("select after insert: %v", err)
	}

	// Sleep a beat so the second now() is observably later, then re-touch.
	time.Sleep(10 * time.Millisecond)
	if err := repo.Touch(t.Context(), u); err != nil {
		t.Fatalf("Touch (update): %v", err)
	}
	var secondSeen time.Time
	if err := pool.QueryRow(t.Context(), "SELECT last_seen_at FROM users WHERE firebase_uid = $1", uid).Scan(&secondSeen); err != nil {
		t.Fatalf("select after update: %v", err)
	}
	if !secondSeen.After(firstSeen) {
		t.Errorf("last_seen_at not bumped: first=%s second=%s", firstSeen, secondSeen)
	}
}

// TestRepo_Touch_PreservesEmailWhenAbsent verifies SPEC-AUTH-019: the COALESCE
// clause keeps the existing email/display_name when a later Touch passes empty
// values (e.g., a token without claims).
func TestRepo_Touch_PreservesEmailWhenAbsent(t *testing.T) {
	pool := testPool(t)
	repo := users.NewRepo(pool)
	uid := "test-uid-coalesce-" + time.Now().Format("150405.000000000")

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = pool.Exec(ctx, "DELETE FROM users WHERE firebase_uid = $1", uid)
	})

	if err := repo.Touch(t.Context(), auth.User{UID: uid, Email: "first@example.com", DisplayName: "First"}); err != nil {
		t.Fatalf("Touch (initial): %v", err)
	}
	if err := repo.Touch(t.Context(), auth.User{UID: uid}); err != nil {
		t.Fatalf("Touch (re-touch with empty fields): %v", err)
	}

	var email, displayName string
	if err := pool.QueryRow(t.Context(), "SELECT email, display_name FROM users WHERE firebase_uid = $1", uid).Scan(&email, &displayName); err != nil {
		t.Fatalf("select: %v", err)
	}
	if email != "first@example.com" {
		t.Errorf("email = %q, want it preserved as %q", email, "first@example.com")
	}
	if displayName != "First" {
		t.Errorf("display_name = %q, want it preserved as %q", displayName, "First")
	}
}
