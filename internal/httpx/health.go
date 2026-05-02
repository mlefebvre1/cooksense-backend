// Package httpx hosts cross-cutting HTTP handlers (health checks, error
// helpers) that are not specific to any domain feature.
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mlefebvre1/cooksense-backend/internal/auth"
)

// Health returns an http.Handler for GET /api/health.
//
// The endpoint is public: no Authorization is required, and no database
// access is performed. Response body: {"status":"ok"}.
//
// SPEC-AUTH-021
func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}

// HealthMe returns an http.Handler for GET /api/health/me. The handler
// expects auth.Middleware to have already populated the request context with
// an authenticated User; if that contract is broken (no user in context) it
// responds with 500 since the misconfiguration is server-side, not a client
// authentication failure.
//
// Response body: {"uid":"…","email":"…"}.
//
// SPEC-AUTH-022
func HealthMe() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := auth.UserFromContext(r.Context())
		if !ok {
			slog.Error("httpx: /api/health/me reached without authenticated user in context")
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]string{
					"code":    "INTERNAL",
					"message": "missing authenticated user",
				},
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"uid":   u.UID,
			"email": u.Email,
		})
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("httpx: encode response", "err", err)
	}
}
