// Command cooksense-server is the HTTP API server for the CookSense application.
//
// Subcommands:
//
//	cooksense-server migrate up        - apply all pending migrations (SPEC-DB-012)
//	cooksense-server migrate down [N]  - roll back N steps, default 1 (SPEC-DB-013)
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
	"github.com/mlefebvre1/cooksense-backend/internal/db"
)

const migrationsDir = "migrations"

func main() {
	if code := run(context.Background(), os.Args[1:]); code != 0 {
		os.Exit(code)
	}
}

func run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Println("cooksense-server starting")
		return 0
	}
	switch args[0] {
	case "migrate":
		return runMigrate(ctx, args[1:])
	case "seed":
		fmt.Println("seed: not implemented")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "cooksense-server: unknown command %q\n", args[0])
		return 1
	}
}

// runMigrate dispatches `migrate up` and `migrate down [N]`.
// SPEC-DB-012, SPEC-DB-013
func runMigrate(ctx context.Context, args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "cooksense-server migrate: missing subcommand (up|down)")
		return 1
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		return 1
	}
	switch args[0] {
	case "up":
		if err := db.Up(ctx, cfg.DatabaseURL, migrationsDir); err != nil {
			slog.Error("migrate up failed", "err", err)
			return 1
		}
		slog.Info("migrations applied")
		return 0
	case "down":
		n := 1
		if len(args) > 1 {
			parsed, err := strconv.Atoi(args[1])
			if err != nil || parsed <= 0 {
				fmt.Fprintf(os.Stderr, "cooksense-server migrate down: N must be a positive integer, got %q\n", args[1])
				return 1
			}
			n = parsed
		}
		if err := db.Down(ctx, cfg.DatabaseURL, migrationsDir, n); err != nil {
			slog.Error("migrate down failed", "err", err)
			return 1
		}
		slog.Info("migrations rolled back", "steps", n)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "cooksense-server migrate: unknown subcommand %q (want up|down)\n", args[0])
		return 1
	}
}
