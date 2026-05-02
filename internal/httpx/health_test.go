package httpx_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mlefebvre1/cooksense-backend/internal/auth"
	"github.com/mlefebvre1/cooksense-backend/internal/httpx"
)

// TestHealthHandler_Public_Returns200 verifies SPEC-AUTH-021: GET /api/health
// is reachable without any Authorization header and returns the documented body.
func TestHealthHandler_Public_Returns200(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	httpx.Health().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("body[status] = %q, want %q", body["status"], "ok")
	}
}

// TestHealthMeHandler_ValidToken_ReturnsUID verifies SPEC-AUTH-022: when the
// request is authenticated, /api/health/me returns the user's uid + email.
func TestHealthMeHandler_ValidToken_ReturnsUID(t *testing.T) {
	const token = "valid"
	want := auth.User{UID: "u-1", Email: "u1@example.com", DisplayName: "U1"}
	v := auth.NewFakeVerifier(map[string]auth.User{token: want})
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mw(httpx.HealthMe()).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["uid"] != want.UID {
		t.Errorf("body[uid] = %q, want %q", body["uid"], want.UID)
	}
	if body["email"] != want.Email {
		t.Errorf("body[email] = %q, want %q", body["email"], want.Email)
	}
}

// TestHealthMeHandler_MissingToken_Returns401 verifies SPEC-AUTH-022 +
// SPEC-AUTH-009: a request without an Authorization header is blocked by the
// middleware before reaching the handler.
func TestHealthMeHandler_MissingToken_Returns401(t *testing.T) {
	v := auth.NewFakeVerifier(nil)
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health/me", nil)
	mw(httpx.HealthMe()).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rr.Code)
	}
	if got := rr.Header().Get("WWW-Authenticate"); got == "" {
		t.Error("WWW-Authenticate header missing on 401")
	}
}

// TestHealthMeHandler_NoUserInContext_Returns500 verifies the defensive path
// in HealthMe: if the handler is reached without auth.Middleware (programming
// error), the response is 500 with the standard envelope rather than a
// confusing zero-valued body.
func TestHealthMeHandler_NoUserInContext_Returns500(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health/me", nil).WithContext(context.Background())
	httpx.HealthMe().ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rr.Code)
	}
}

type noopToucher struct{}

func (noopToucher) Touch(_ context.Context, _ auth.User) error { return nil }
