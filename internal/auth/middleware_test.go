package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mlefebvre1/cooksense-backend/internal/auth"
)

const (
	validToken = "valid-token"
	wwwAuthExp = `Bearer realm="cooksense"`
)

func validUser() auth.User {
	return auth.User{UID: "uid-123", Email: "alice@example.com", DisplayName: "Alice"}
}

type noopToucher struct{}

func (noopToucher) Touch(_ context.Context, _ auth.User) error { return nil }

type errorToucher struct{ err error }

func (t errorToucher) Touch(_ context.Context, _ auth.User) error { return t.err }

type spyToucher struct {
	called bool
	last   auth.User
}

func (s *spyToucher) Touch(_ context.Context, u auth.User) error {
	s.called = true
	s.last = u
	return nil
}

// okHandler asserts that the request reached the downstream handler with the
// expected user in context, and writes a 200.
func okHandler(t *testing.T, want auth.User) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := auth.UserFromContext(r.Context())
		if !ok {
			t.Errorf("UserFromContext = (_, false), want true")
			http.Error(w, "no user", http.StatusInternalServerError)
			return
		}
		if got != want {
			t.Errorf("UserFromContext = %+v, want %+v", got, want)
		}
		w.WriteHeader(http.StatusOK)
	})
}

func notReachedHandler(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Error("downstream handler should not be reached on a 401")
	})
}

func assertUnauthnEnvelope(t *testing.T, rr *httptest.ResponseRecorder) {
	t.Helper()
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rr.Code)
	}
	if got := rr.Header().Get("WWW-Authenticate"); got != wwwAuthExp {
		t.Errorf("WWW-Authenticate = %q, want %q", got, wwwAuthExp)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v (raw=%q)", err, rr.Body.String())
	}
	if body.Error.Code != "UNAUTHENTICATED" {
		t.Errorf("error.code = %q, want UNAUTHENTICATED", body.Error.Code)
	}
	if body.Error.Message == "" {
		t.Error("error.message must not be empty")
	}
}

// TestMiddleware_MissingAuthorizationHeader_Returns401 verifies SPEC-AUTH-009,
// SPEC-AUTH-012, SPEC-AUTH-023.
func TestMiddleware_MissingAuthorizationHeader_Returns401(t *testing.T) {
	v := auth.NewFakeVerifier(map[string]auth.User{validToken: validUser()})
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	mw(notReachedHandler(t)).ServeHTTP(rr, req)

	assertUnauthnEnvelope(t, rr)
}

// TestMiddleware_MalformedAuthorizationHeader_Returns401 verifies
// SPEC-AUTH-010, SPEC-AUTH-012, SPEC-AUTH-023.
func TestMiddleware_MalformedAuthorizationHeader_Returns401(t *testing.T) {
	cases := []string{
		"Basic abc",    // wrong scheme
		"Bearer",       // no token
		"Bearer ",      // empty token after prefix
		"bearer token", // lowercase scheme not accepted
		"Token abc",    // unsupported scheme
	}
	for _, header := range cases {
		t.Run(header, func(t *testing.T) {
			v := auth.NewFakeVerifier(map[string]auth.User{validToken: validUser()})
			mw := auth.Middleware(v, noopToucher{})

			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set("Authorization", header)
			mw(notReachedHandler(t)).ServeHTTP(rr, req)

			assertUnauthnEnvelope(t, rr)
		})
	}
}

// TestMiddleware_ExpiredOrInvalidToken_Returns401 verifies SPEC-AUTH-011,
// SPEC-AUTH-012.
func TestMiddleware_ExpiredOrInvalidToken_Returns401(t *testing.T) {
	v := auth.NewFakeVerifier(map[string]auth.User{validToken: validUser()})
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer expired-or-unknown")
	mw(notReachedHandler(t)).ServeHTTP(rr, req)

	assertUnauthnEnvelope(t, rr)
}

// TestMiddleware_WrongAudienceToken_Returns401 verifies SPEC-AUTH-011: any
// non-nil error from the verifier is mapped to a 401.
func TestMiddleware_WrongAudienceToken_Returns401(t *testing.T) {
	v := &errVerifier{err: errors.New("wrong audience")}
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	mw(notReachedHandler(t)).ServeHTTP(rr, req)

	assertUnauthnEnvelope(t, rr)
}

// TestMiddleware_ValidToken_CallsNextAndSetsContext verifies SPEC-AUTH-013,
// SPEC-AUTH-014, SPEC-AUTH-016.
func TestMiddleware_ValidToken_CallsNextAndSetsContext(t *testing.T) {
	want := validUser()
	v := auth.NewFakeVerifier(map[string]auth.User{validToken: want})
	spy := &spyToucher{}
	mw := auth.Middleware(v, spy)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	mw(okHandler(t, want)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}
	if !spy.called {
		t.Error("Toucher.Touch was not called on successful auth")
	}
	if spy.last != want {
		t.Errorf("Touch received user %+v, want %+v", spy.last, want)
	}
}

// TestMiddleware_TouchError_ContinuesRequest verifies SPEC-AUTH-020: a Touch
// failure is non-fatal — the request continues and reaches the handler.
func TestMiddleware_TouchError_ContinuesRequest(t *testing.T) {
	want := validUser()
	v := auth.NewFakeVerifier(map[string]auth.User{validToken: want})
	mw := auth.Middleware(v, errorToucher{err: errors.New("db unavailable")})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	mw(okHandler(t, want)).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 despite Touch error", rr.Code)
	}
}

// TestMiddleware_DoesNotReadBody verifies SPEC-AUTH-015: the middleware short-
// circuits a 401 without consuming r.Body.
func TestMiddleware_DoesNotReadBody(t *testing.T) {
	body := &countingReader{src: strings.NewReader("payload")}
	v := auth.NewFakeVerifier(nil)
	mw := auth.Middleware(v, noopToucher{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/protected", body)
	// No Authorization header — middleware should reject before reading body.
	mw(notReachedHandler(t)).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rr.Code)
	}
	if body.reads > 0 {
		t.Errorf("body Read called %d times, want 0", body.reads)
	}
}

// TestUserFromContext_NoUser_ReturnsFalse verifies SPEC-AUTH-016: the accessor
// reports false when no user is present.
func TestUserFromContext_NoUser_ReturnsFalse(t *testing.T) {
	if u, ok := auth.UserFromContext(context.Background()); ok {
		t.Errorf("UserFromContext on empty ctx = (%+v, true), want (User{}, false)", u)
	}
	if u, ok := auth.UserFromContext(nil); ok { //nolint:staticcheck // SA1012: explicitly testing nil-context safety
		t.Errorf("UserFromContext on nil ctx = (%+v, true), want (User{}, false)", u)
	}
}

// TestFakeVerifier_UnknownToken_Errors verifies SPEC-AUTH-006.
func TestFakeVerifier_UnknownToken_Errors(t *testing.T) {
	f := auth.NewFakeVerifier(map[string]auth.User{validToken: validUser()})
	if _, err := f.Verify(t.Context(), "nope"); err == nil {
		t.Error("Verify(unknown) returned nil error, want error")
	}
	got, err := f.Verify(t.Context(), validToken)
	if err != nil {
		t.Fatalf("Verify(valid) error: %v", err)
	}
	if got != validUser() {
		t.Errorf("Verify(valid) = %+v, want %+v", got, validUser())
	}
}

// errVerifier is a Verifier that always returns a fixed error, used to model
// "wrong audience" and similar SDK failures (SPEC-AUTH-011).
type errVerifier struct{ err error }

func (e *errVerifier) Verify(_ context.Context, _ string) (auth.User, error) {
	return auth.User{}, e.err
}

// countingReader counts how many times Read is invoked, so a test can assert
// the middleware never touched the request body (SPEC-AUTH-015).
type countingReader struct {
	src   io.Reader
	reads int
}

func (c *countingReader) Read(p []byte) (int, error) {
	c.reads++
	return c.src.Read(p)
}
