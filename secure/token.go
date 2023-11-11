package secure

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	TOKEN_TYPE_BEARER  = "Bearer"
	TOKEN_TYPE_REFRESH = "Refresh"
)

// Token is used to authenticate a user.
// All fields are private so will not be modified outside of this package.
type Token struct {
	id        string
	ttype     string // token type
	realm     string
	client    string
	subject   string
	scope     []string
	issuedAt  time.Time
	expiresAt time.Time
}

// NewToken creates a new token.
func NewToken(ttype, realm, client, subject string, scope []string) *Token {
	return &Token{
		ttype:   ttype,
		realm:   realm,
		client:  client,
		subject: subject,
		scope:   scope,
	}
}

func (t *Token) Realm() string {
	return t.realm
}

func (t *Token) Client() string {
	return t.client
}

func (t *Token) Subject() string {
	return t.subject
}

func (t *Token) Scope() []string {
	return t.scope
}

func (t *Token) HasScope(scope string) bool {
	if t != nil {
		for _, s := range t.scope {
			if strings.EqualFold(s, scope) {
				return true
			}
		}
	}
	return false
}

func (t *Token) IssuedAt() time.Time {
	return t.issuedAt
}

func (t *Token) ExpiresAt() time.Time {
	return t.expiresAt
}

func (t *Token) IsExpired() bool {
	return !t.expiresAt.IsZero() && t.expiresAt.After(time.Now().UTC())
}

// TokenStore used to manage tokens.
type TokenStore interface {
	// Issue issues a new token with the given ttl.
	Issue(token *Token, ttl time.Duration) (string, error)
	// Renew renews the token and returns the new one.
	Renew(value string) (string, error)
	// Verify verifies the token and returns the token if valid.
	Verify(value string) (*Token, error)
	// Revoke revokes the token and returns the token if revoked.
	Revoke(value string) (*Token, error)
}

// InMemoryTokenStore is an in-memory token store.
type InMemoryTokenStore struct {
	tokens sync.Map
}

var _ TokenStore = (*InMemoryTokenStore)(nil)

func (s *InMemoryTokenStore) Issue(token *Token, ttl time.Duration) (string, error) {
	token.id = uuid.NewString()
	token.issuedAt = time.Now().UTC()
	token.expiresAt = token.issuedAt.Add(ttl)
	s.tokens.Store(token.id, token)
	return token.id, nil
}

func (s *InMemoryTokenStore) Renew(value string) (string, error) {
	if value, exists := s.tokens.Load(value); exists {
		token := value.(*Token)
		token.issuedAt = time.Now().UTC()
		token.expiresAt = token.issuedAt.Add(token.expiresAt.Sub(token.issuedAt))
		return token.id, nil
	}
	return "", ErrInvalidToken
}

func (s *InMemoryTokenStore) Verify(value string) (*Token, error) {
	if value, exists := s.tokens.Load(value); exists {
		token := value.(*Token)
		if token.IsExpired() {
			s.tokens.Delete(token.id)
			return token, ErrInvalidToken
		}
		return token, nil
	}
	return nil, ErrInvalidToken
}

func (s *InMemoryTokenStore) Revoke(value string) (*Token, error) {
	if value, exists := s.tokens.LoadAndDelete(value); exists {
		token := value.(*Token)
		return token, nil
	}
	return nil, nil
}
