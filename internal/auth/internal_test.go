package auth

import (
	"context"
	"testing"

	firebaseauth "firebase.google.com/go/v4/auth"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
)

// TestUserFromToken_ExtractsClaims verifies SPEC-AUTH-005: firebaseVerifier
// derives User.UID, Email and DisplayName from the decoded token's claims.
func TestUserFromToken_ExtractsClaims(t *testing.T) {
	tok := &firebaseauth.Token{
		UID: "uid-9",
		Claims: map[string]any{
			"email": "alice@example.com",
			"name":  "Alice",
		},
	}
	got := userFromToken(tok)
	want := User{UID: "uid-9", Email: "alice@example.com", DisplayName: "Alice"}
	if got != want {
		t.Errorf("userFromToken = %+v, want %+v", got, want)
	}
}

// TestUserFromToken_MissingClaims_ReturnsZeroFields verifies SPEC-AUTH-005:
// missing or wrongly-typed claims are tolerated; only UID is mandatory.
func TestUserFromToken_MissingClaims_ReturnsZeroFields(t *testing.T) {
	tok := &firebaseauth.Token{UID: "uid-only"}
	got := userFromToken(tok)
	want := User{UID: "uid-only"}
	if got != want {
		t.Errorf("userFromToken (no claims) = %+v, want %+v", got, want)
	}
}

// TestNewFirebaseApp_NilConfig_ReturnsError verifies SPEC-AUTH-001 input
// validation: a nil config is rejected without attempting network init.
func TestNewFirebaseApp_NilConfig_ReturnsError(t *testing.T) {
	if _, err := NewFirebaseApp(context.Background(), nil); err == nil {
		t.Error("NewFirebaseApp(nil) returned nil error, want error")
	}
}

// TestNewFirebaseApp_ValidConfig_BuildsAppLazily verifies SPEC-AUTH-001: the
// constructor accepts a fully-populated config and returns a non-nil app even
// if the credentials file is not yet present (the SDK validates lazily on
// first use). The check guards us against future regressions where the
// constructor starts eagerly opening the file at the wrong layer.
func TestNewFirebaseApp_ValidConfig_BuildsAppLazily(t *testing.T) {
	cfg := &config.Config{
		FirebaseProjectID:    "test-project",
		GoogleAppCredentials: "/does/not/exist.json",
	}
	app, err := NewFirebaseApp(context.Background(), cfg)
	if err != nil {
		t.Fatalf("NewFirebaseApp returned error: %v", err)
	}
	if app == nil {
		t.Fatal("NewFirebaseApp returned nil app")
	}
}

// TestNewFirebaseVerifier_NilApp_ReturnsError verifies SPEC-AUTH-005 input
// validation.
func TestNewFirebaseVerifier_NilApp_ReturnsError(t *testing.T) {
	if _, err := NewFirebaseVerifier(nil); err == nil {
		t.Error("NewFirebaseVerifier(nil) returned nil error, want error")
	}
}

type stubIDTokenVerifier struct {
	tok *firebaseauth.Token
	err error
}

func (s stubIDTokenVerifier) VerifyIDToken(_ context.Context, _ string) (*firebaseauth.Token, error) {
	return s.tok, s.err
}

// TestFirebaseVerifier_Verify_PropagatesSDKError verifies SPEC-AUTH-005: any
// error returned by VerifyIDToken is forwarded as-is so the middleware can
// translate it to 401 (SPEC-AUTH-011).
func TestFirebaseVerifier_Verify_PropagatesSDKError(t *testing.T) {
	v := &firebaseVerifier{client: stubIDTokenVerifier{err: errFake("boom")}}
	if _, err := v.Verify(context.Background(), "tok"); err == nil {
		t.Error("Verify with stub error returned nil, want propagated error")
	}
}

// TestFirebaseVerifier_Verify_MapsToken verifies SPEC-AUTH-005: a successful
// VerifyIDToken response is mapped through userFromToken.
func TestFirebaseVerifier_Verify_MapsToken(t *testing.T) {
	tok := &firebaseauth.Token{
		UID:    "uid-x",
		Claims: map[string]any{"email": "x@example.com", "name": "X"},
	}
	v := &firebaseVerifier{client: stubIDTokenVerifier{tok: tok}}
	got, err := v.Verify(context.Background(), "tok")
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	want := User{UID: "uid-x", Email: "x@example.com", DisplayName: "X"}
	if got != want {
		t.Errorf("Verify = %+v, want %+v", got, want)
	}
}

type errFake string

func (e errFake) Error() string { return string(e) }
