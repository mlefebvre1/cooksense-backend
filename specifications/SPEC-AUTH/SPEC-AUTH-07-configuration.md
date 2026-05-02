# SPEC-AUTH-07 — Configuration

← [Index](SPEC-AUTH-00-index.md)

---

## 1. Environment variables

**SPEC-AUTH-002** and **SPEC-AUTH-003** require two new variables to be added
to the existing `config.Config` struct and validated at startup.

| Variable | `Config` field | Required | Validation rule |
|----------|---------------|----------|-----------------|
| `FIREBASE_PROJECT_ID` | `FirebaseProjectID string` | yes | Non-empty string |
| `GOOGLE_APPLICATION_CREDENTIALS` | `GoogleAppCredentials string` | yes | Non-empty string (path; existence not validated at config load time — the Firebase SDK will fail on init if the file is missing) |

Both variables are already listed in `.env.example` (verified). If either is
missing or empty, `config.Load` shall return a descriptive error and the
server shall exit non-zero. No default values shall be supplied in code.

## 2. Credential file rules

**SPEC-AUTH-024** — A placeholder file `secrets/firebase-admin.json.example`
shall be committed to the repository. It shall contain a JSON skeleton with
no real values:

```json
{
  "type": "service_account",
  "project_id": "YOUR_FIREBASE_PROJECT_ID",
  "private_key_id": "YOUR_KEY_ID",
  "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-XXXX@YOUR_PROJECT_ID.iam.gserviceaccount.com",
  "client_id": "YOUR_CLIENT_ID",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token"
}
```

**SPEC-AUTH-025** — `.env.example` shall contain (already present, confirmed):

```
FIREBASE_PROJECT_ID=changeme
GOOGLE_APPLICATION_CREDENTIALS=path/to/firebase-admin.json
```

The `.gitignore` rule `secrets/*.json` (Story 01, SPEC-BOOT) ensures the real
credential file is never committed. Only `secrets/*.json.example` files are
tracked.

## 3. Loading implementation

`config.Load` (already introduced in Story 03 as `internal/config/config.go`)
shall be extended with the two new fields. The loading strategy remains
`os.Getenv` with explicit required-field validation. No third-party config
library shall be introduced.

```go
// In internal/config/config.go — additions for Story 04
type Config struct {
    // ... existing fields (AppPort, LogLevel, LogFormat, DatabaseURL) ...

    // SPEC-AUTH-002
    FirebaseProjectID string

    // SPEC-AUTH-003
    GoogleAppCredentials string
}
```
