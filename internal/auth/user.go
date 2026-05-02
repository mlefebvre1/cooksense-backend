package auth

// User represents an authenticated Firebase principal.
//
// SPEC-AUTH-007
type User struct {
	UID         string
	Email       string
	DisplayName string
}

// userKey is the unexported context key under which an authenticated User is
// stored by Middleware and retrieved by UserFromContext.
//
// SPEC-AUTH-014, SPEC-AUTH-016
type userKey struct{}
