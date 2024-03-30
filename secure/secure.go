package secure

import (
	"errors"
	"fmt"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/rueidis"
)

const (
	AUTH_HEADER_KEY    = "Authorization"
	AUTH_SCHEMA_BASIC  = "basic"
	AUTH_SCHEMA_BEARER = "bearer"
)

// NewTokenStore returns a new TokenStore with the given config.
// There are three types of token stores: jwt, redis, and memory.
// rueidis.Client is required if the type of token store is redis.
func NewTokenStore(cfg config.SecureTokenConfig, rdb rueidis.Client) (TokenStore, error) {
	switch cfg.GetStore() {
	case "jwt":
		alg := jwt.GetSigningMethod(cfg.GetSigningMethod())
		if alg == nil {
			return nil, fmt.Errorf("unknown signing method: %s", cfg.GetSigningMethod())
		}
		return NewJsonWebTokenStore(cfg.GetIssuer(), cfg.GetAudience(), cfg.GetSigningMethod(), cfg.GetPrivateKey(), cfg.GetPublicKey())
	case "redis":
		if rdb == nil {
			return nil, errors.New("redis client is required for redis token store")
		}
		return NewRedisTokenStore(rdb, cfg.GetBucket())
	case "memory":
		return &InMemoryTokenStore{}, nil
	}
	return nil, fmt.Errorf("unknown token store: %s", cfg.GetStore())
}
