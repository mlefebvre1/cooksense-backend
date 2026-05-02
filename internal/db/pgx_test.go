//go:build integration

package db_test

import (
	"testing"
)

// TestOpen_WithRealPostgres_ReturnsPool verifies SPEC-DB-004, SPEC-DB-005, SPEC-DB-006.
func TestOpen_WithRealPostgres_ReturnsPool(t *testing.T) {
	pool := testDB(t)
	if pool == nil {
		t.Fatal("expected non-nil pool")
	}
	if err := pool.Ping(t.Context()); err != nil {
		t.Fatalf("pool.Ping: %v", err)
	}
	cfg := pool.Config()
	if cfg.MinConns != 2 {
		t.Errorf("MinConns = %d, want 2", cfg.MinConns)
	}
	if cfg.MaxConns != 10 {
		t.Errorf("MaxConns = %d, want 10", cfg.MaxConns)
	}
}
