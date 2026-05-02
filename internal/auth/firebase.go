package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"github.com/mlefebvre1/cooksense-backend/internal/config"
)

// NewFirebaseApp initializes the Firebase Admin SDK from cfg. It is intended
// to be called exactly once at server startup; the returned *firebase.App is
// shared via dependency injection.
//
// SPEC-AUTH-001, SPEC-AUTH-002, SPEC-AUTH-003
func NewFirebaseApp(ctx context.Context, cfg *config.Config) (*firebase.App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("auth: nil config")
	}
	fbCfg := &firebase.Config{ProjectID: cfg.FirebaseProjectID}
	opt := option.WithCredentialsFile(cfg.GoogleAppCredentials)
	app, err := firebase.NewApp(ctx, fbCfg, opt)
	if err != nil {
		return nil, fmt.Errorf("auth: init firebase app: %w", err)
	}
	return app, nil
}

// NewFirebaseVerifier creates a Verifier backed by the given Firebase app.
//
// SPEC-AUTH-005
func NewFirebaseVerifier(app *firebase.App) (Verifier, error) {
	if app == nil {
		return nil, fmt.Errorf("auth: nil firebase app")
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("auth: firebase auth client: %w", err)
	}
	return &firebaseVerifier{client: client}, nil
}
