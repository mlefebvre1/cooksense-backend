package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlefebvre1/cooksense-backend/internal/auth"
)

// Toucher upserts the user row and updates last_seen_at on every
// authenticated request, providing lazy provisioning.
//
// SPEC-AUTH-017
type Toucher interface {
	Touch(ctx context.Context, u auth.User) error
}

// Repo implements Toucher and other user-related persistence backed by a
// shared pgx pool injected at construction time.
//
// SPEC-AUTH-018
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo returns a new Repo backed by pool.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// touchSQL is the canonical UPSERT statement defined by SPEC-AUTH-019.
const touchSQL = `INSERT INTO users (firebase_uid, display_name, email)
VALUES ($1, $2, $3)
ON CONFLICT (firebase_uid) DO UPDATE
   SET last_seen_at  = now(),
       email         = COALESCE(EXCLUDED.email, users.email),
       display_name  = COALESCE(EXCLUDED.display_name, users.display_name);`

// Touch inserts the user row keyed by Firebase UID, or updates last_seen_at
// (and any newly-supplied email/display_name) on conflict.
//
// SPEC-AUTH-018, SPEC-AUTH-019
func (r *Repo) Touch(ctx context.Context, u auth.User) error {
	var displayName, email any
	if u.DisplayName != "" {
		displayName = u.DisplayName
	}
	if u.Email != "" {
		email = u.Email
	}
	if _, err := r.pool.Exec(ctx, touchSQL, u.UID, displayName, email); err != nil {
		return fmt.Errorf("users: touch %q: %w", u.UID, err)
	}
	return nil
}
