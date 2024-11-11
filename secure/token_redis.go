package secure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

var (
	luaIssue = rueidis.NewLuaScript(`
        local keys = redis.call('keys', ARGV[1])
        for _, key in ipairs(keys) do
          redis.call('del', key)
        end
		redis.call('set', table.concat({KEY[1], '.ttl'}), ARGV[3], 'EX', ARGV[3])
        return redis.call('set', KEY[1], ARGV[2], 'EX', ARGV[3])
    `)
)

// RedisTokenStore is a token store that uses Redis to store tokens.
type RedisTokenStore struct {
	rdb rueidis.Client
	bkt string
}

var _ TokenStore = (*RedisTokenStore)(nil)

type RedisTokenConfig interface {
	GetStore() string
	GetBucket() string
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
	GetIssuer() string
	GetAudience() string
	GetSigningMethod() string
	GetPublicKey() []byte
	GetPrivateKey() []byte
}

// NewRedisTokenStore creates a new Redis token store.
func NewRedisTokenStore(rdb rueidis.Client, bkt string) (*RedisTokenStore, error) {
	return &RedisTokenStore{
		rdb: rdb,
		bkt: bkt,
	}, nil
}

func (s *RedisTokenStore) issue(ctx context.Context, token *Token, ttl time.Duration) (string, error) {
	token.issuedAt = time.Now().UTC()
	token.expiresAt = token.issuedAt.Add(ttl)
	td, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf("%s:%s", s.bkt, token.id)
	var rr rueidis.RedisResult
	if token.ttype == TokenTypeBearer {
		rr = luaIssue.Exec(ctx, s.rdb, []string{
			key, // KEY[1]: token key
		}, []string{
			fmt.Sprintf("%s:%s:*", s.bkt, token.subject), // ARGV[1]: search pattern
			string(td), // ARGV[2]: token data
			strconv.FormatInt(int64(ttl.Seconds()), 10), // ARGV[3]: expiration
		})
	} else {
		return "", ErrUnsupportedTokenType
	}
	if err := rr.Error(); err != nil {
		return "", err
	}
	return token.id, nil
}

func (s *RedisTokenStore) verify(ctx context.Context, value string, del bool) (*Token, error) {
	var td string
	var err error
	var cmd rueidis.Completed
	if del {
		cmd = s.rdb.B().Getdel().Key(fmt.Sprintf("%s:%s", s.bkt, value)).Build()
	} else {
		cmd = s.rdb.B().Get().Key(fmt.Sprintf("%s:%s", s.bkt, value)).Build()
	}
	td, err = s.rdb.Do(ctx, cmd).ToString()
	if td == "" || errors.Is(err, rueidis.Nil) {
		return nil, ErrInvalidToken
	} else if err != nil {
		return nil, err
	}
	token := &Token{}
	if err := json.Unmarshal([]byte(td), token); err != nil {
		return nil, err
	}
	token.id = value
	return token, nil
}

func (s *RedisTokenStore) Issue(ctx context.Context, token *Token, ttl time.Duration) (string, error) {
	token.id = fmt.Sprintf("%s:%s:%s", token.client, token.subject, uuid.NewString())
	return s.issue(ctx, token, ttl)
}

func (s *RedisTokenStore) Renew(ctx context.Context, value string, ttl time.Duration) (string, error) {
	token, err := s.verify(ctx, value, false)
	if err != nil {
		return "", err
	}
	return s.issue(ctx, token, ttl)
}

func (s *RedisTokenStore) Verify(ctx context.Context, value string) (*Token, error) {
	token, err := s.verify(ctx, value, false)
	if err != nil {
		return nil, err
	}
	if _, err := s.issue(ctx, token, token.expiresAt.Sub(token.issuedAt)); err != nil {
		return nil, err
	}
	return token, nil
}

func (s *RedisTokenStore) Revoke(ctx context.Context, value string) (*Token, error) {
	return s.verify(ctx, value, true)
}
