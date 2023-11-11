package secure

import (
	"context"
)

type identityKey struct{}

// IdentityFromContext returns the identity from the given context.
func IdentityFromContext(ctx context.Context) *Identity {
	if t, b := ctx.Value(identityKey{}).(*Identity); b {
		return t
	} else {
		return nil
	}
}
