package secure

import (
	"context"
)

type identityKey struct{}

// IdentityFromContext returns the identity from the given context.
func IdentityFromContext(ctx context.Context) *Identity {
	if id, ok := ctx.Value(identityKey{}).(*Identity); ok {
		return id
	}
	return nil
}

// SubjectFromContext returns the subject of the token from the given context.
func SubjectFromContext(ctx context.Context) (string, bool) {
	if id := IdentityFromContext(ctx); id != nil && id.token != nil {
		return id.token.subject, true
	}
	return "", false
}
