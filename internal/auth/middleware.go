package auth

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// Toucher is the consumer-side interface that Middleware uses to record a
// per-request "user seen" event. Implementations live in internal/users; auth
// declares the structural type locally so the auth package does not import
// the users package (Clean-Architecture dependency direction, see SPEC-AUTH §5).
//
// Any concrete type whose method set matches this interface (notably
// *users.Repo, SPEC-AUTH-018) satisfies it.
//
// SPEC-AUTH-008, SPEC-AUTH-013, SPEC-AUTH-017
type Toucher interface {
	Touch(ctx context.Context, u User) error
}

const (
	authHeader       = "Authorization"
	bearerPrefix     = "Bearer "
	wwwAuthenticate  = "WWW-Authenticate"
	wwwAuthValue     = `Bearer realm="cooksense"`
	contentType      = "Content-Type"
	contentTypeJSON  = "application/json"
	codeUnauthn      = "UNAUTHENTICATED"
	msgMissingHeader = "missing or invalid Authorization header"
	msgInvalidToken  = "invalid or expired token"
)

// Middleware returns an http.Handler wrapper that verifies a Firebase ID
// token on every request, lazily provisions the matching local user row,
// and injects the authenticated User into the request context.
//
// SPEC-AUTH-008, SPEC-AUTH-009, SPEC-AUTH-010, SPEC-AUTH-011,
// SPEC-AUTH-012, SPEC-AUTH-013, SPEC-AUTH-014, SPEC-AUTH-015
func Middleware(v Verifier, toucher Toucher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractBearerToken(r)
			if !ok {
				writeUnauthn(w, msgMissingHeader)
				return
			}

			user, err := v.Verify(r.Context(), token)
			if err != nil {
				slog.Debug("auth: verify failed", "err", err)
				writeUnauthn(w, msgInvalidToken)
				return
			}

			if toucher != nil {
				if terr := toucher.Touch(r.Context(), user); terr != nil {
					slog.Warn("auth: touch failed", "firebase_uid", user.UID, "err", terr)
				}
			}

			ctx := context.WithValue(r.Context(), userKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractBearerToken returns the Bearer token from the Authorization header,
// or false if the header is absent or malformed.
//
// SPEC-AUTH-009, SPEC-AUTH-010
func extractBearerToken(r *http.Request) (string, bool) {
	h := r.Header.Get(authHeader)
	if h == "" {
		return "", false
	}
	if !strings.HasPrefix(h, bearerPrefix) {
		return "", false
	}
	token := strings.TrimPrefix(h, bearerPrefix)
	if token == "" {
		return "", false
	}
	return token, true
}

// writeUnauthn writes the standard 401 envelope with the WWW-Authenticate
// challenge header.
//
// SPEC-AUTH-012, SPEC-AUTH-023
func writeUnauthn(w http.ResponseWriter, message string) {
	w.Header().Set(wwwAuthenticate, wwwAuthValue)
	w.Header().Set(contentType, contentTypeJSON)
	w.WriteHeader(http.StatusUnauthorized)
	body := errorEnvelope{Error: errorBody{Code: codeUnauthn, Message: message}}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("auth: write 401 body", "err", err)
	}
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// UserFromContext retrieves the authenticated User from ctx. The boolean is
// false when no user is present (handler reached without going through
// Middleware, or middleware bypass).
//
// SPEC-AUTH-016
func UserFromContext(ctx context.Context) (User, bool) {
	if ctx == nil {
		return User{}, false
	}
	u, ok := ctx.Value(userKey{}).(User)
	if !ok {
		return User{}, false
	}
	return u, true
}
