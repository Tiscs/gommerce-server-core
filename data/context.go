package data

import (
	"context"
)

type idWorkderKey struct{}

func IdWorkerFromContext(ctx context.Context) IdWorker {
	if idw, ok := ctx.Value(idWorkderKey{}).(IdWorker); ok {
		return idw
	}
	return nil
}
