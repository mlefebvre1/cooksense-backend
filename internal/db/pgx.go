// Package db provides the PostgreSQL connection pool and migration runner for cooksense-server.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
)

const (
	defaultMinConns          = 2
	defaultMaxConns          = 10
	defaultHealthCheckPeriod = 1 * time.Minute
)

// Open creates and validates a PostgreSQL connection pool using cfg.DatabaseURL.
// It applies pool defaults (min 2, max 10, healthcheck 1 minute) and pings the
// database before returning. The variadic overrides are applied after the
// defaults and before the pool is created, allowing tests to tune the config
// without changing the production call site.
//
// Returns an error if the DSN cannot be parsed or the database is unreachable.
//
// SPEC-DB-004, SPEC-DB-005, SPEC-DB-006, SPEC-DB-007
func Open(ctx context.Context, cfg config.Config, overrides ...func(*pgxpool.Config)) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: parse DATABASE_URL: %w", err)
	}
	poolCfg.MinConns = defaultMinConns
	poolCfg.MaxConns = defaultMaxConns
	poolCfg.HealthCheckPeriod = defaultHealthCheckPeriod
	for _, override := range overrides {
		override(poolCfg)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}
	return pool, nil
}
