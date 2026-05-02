package auth

import (
	"context"
	"fmt"

	firebaseauth "firebase.google.com/go/v4/auth"
)

// Verifier verifies a Firebase ID token and returns the caller's identity.
//
// SPEC-AUTH-004
type Verifier interface {
	Verify(ctx context.Context, idToken string) (User, error)
}

// idTokenVerifier captures the slice of *firebaseauth.Client we depend on.
// It exists so firebaseVerifier can be unit-tested without a real Firebase
// project; production code always passes a real *firebaseauth.Client.
type idTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*firebaseauth.Token, error)
}

// firebaseVerifier implements Verifier using the Firebase Admin SDK.
//
// SPEC-AUTH-005
type firebaseVerifier struct {
	client idTokenVerifier
}

// Verify decodes idToken via the Firebase Admin SDK and returns the
// corresponding User. Any SDK error is returned as-is; the middleware maps it
// to a 401 response.
//
// SPEC-AUTH-005
func (v *firebaseVerifier) Verify(ctx context.Context, idToken string) (User, error) {
	tok, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return User{}, err
	}
	return userFromToken(tok), nil
}

// userFromToken extracts a User value from a verified Firebase token.
// Email and display name come from the claims map populated by the SDK.
func userFromToken(tok *firebaseauth.Token) User {
	u := User{UID: tok.UID}
	if email, ok := tok.Claims["email"].(string); ok {
		u.Email = email
	}
	if name, ok := tok.Claims["name"].(string); ok {
		u.DisplayName = name
	}
	return u
}

// FakeVerifier is an in-memory Verifier for use in tests only. It maps tokens
// to known User values; unknown tokens yield an "unknown token" error.
//
// SPEC-AUTH-006
type FakeVerifier struct {
	tokens map[string]User
}

// NewFakeVerifier returns a FakeVerifier seeded with the given token-to-user
// mapping. The map is consulted by Verify; callers retain ownership.
//
// SPEC-AUTH-006
func NewFakeVerifier(tokens map[string]User) *FakeVerifier {
	return &FakeVerifier{tokens: tokens}
}

// Verify returns the User registered for idToken, or an error for unknown
// tokens. No network call is made.
//
// SPEC-AUTH-006
func (f *FakeVerifier) Verify(_ context.Context, idToken string) (User, error) {
	if u, ok := f.tokens[idToken]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("unknown token")
}
