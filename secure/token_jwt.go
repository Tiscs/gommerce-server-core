package secure

import (
	"context"
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JsonWebTokenStore is a token store that uses JSON Web Tokens (JWT) to store tokens.
type JsonWebTokenStore struct {
	tokenIssuer   string
	tokenAudience string
	signingMethod jwt.SigningMethod
	signingKey    any
	verifyKey     any
}

var _ TokenStore = (*JsonWebTokenStore)(nil)

// ExtendedClaims is a custom claims type that extends the default claims with additional claims.
type ExtendedClaims struct {
	jwt.RegisteredClaims
	// the `typ` (Type) claim. A custom claim to identify the type of the token.
	Type string `json:"typ,omitempty"`
	// the `realm` (Realm) claim. A custom claim to identify the realm of the token.
	Realm string `json:"realm,omitempty"`
	// the `azp` (Authorized party) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
	Client string `json:"azp,omitempty"`
	// the `scope` (Scope) claim. See https://datatracker.ietf.org/doc/html/rfc6749#section-3.3
	// Note: the scope claim is a space-separated list of scopes, not a JSON array.
	Scope string `json:"scope,omitempty"`
}

type keyData []byte

func (k keyData) toRSAPrivateKey() (*rsa.PrivateKey, error) {
	if k == nil {
		return nil, nil
	}
	return jwt.ParseRSAPrivateKeyFromPEM(k)
}

func (k keyData) toRSAPublicKey() (*rsa.PublicKey, error) {
	if k == nil {
		return nil, nil
	}
	return jwt.ParseRSAPublicKeyFromPEM(k)
}

// NewJsonWebTokenStore creates a new JSON Web Token (JWT) token store.
func NewJsonWebTokenStore(tokenIssuer, tokenAudience, signingAlgName string, signingKeyData, verifyKeyData keyData) (*JsonWebTokenStore, error) {
	// jwt.MarshalSingleStringAsArray is a global variable that controls whether a single string value is marshalled as a string or an array of strings.
	// set to true to marshal single string values as arrays
	jwt.MarshalSingleStringAsArray = false
	var signingKey any
	var verifyKey any
	signingMethod := jwt.GetSigningMethod(signingAlgName)
	if signingMethod == jwt.SigningMethodRS256 || signingMethod == jwt.SigningMethodRS384 || signingMethod == jwt.SigningMethodRS512 {
		// if the signing method is RSA, parse the private key from the signing key data
		if privateKey, err := signingKeyData.toRSAPrivateKey(); err != nil {
			return nil, err
		} else {
			signingKey = privateKey
		}
		if publicKey, err := verifyKeyData.toRSAPublicKey(); err != nil {
			return nil, err
		} else {
			verifyKey = publicKey
		}
	} else if signingMethod == jwt.SigningMethodHS256 || signingMethod == jwt.SigningMethodHS384 || signingMethod == jwt.SigningMethodHS512 {
		// if the signing method is HMAC, use the signing key data as the signing key and verify key
		signingKey = signingKeyData
		verifyKey = verifyKeyData
	} else {
		// else, return an error
		return nil, ErrUnsupportedSigningMethod
	}
	// return a new JSON Web Token (JWT) token store
	return &JsonWebTokenStore{
		tokenIssuer:   tokenIssuer,
		tokenAudience: tokenAudience,
		signingMethod: signingMethod,
		signingKey:    signingKey,
		verifyKey:     verifyKey,
	}, nil
}

func (s *JsonWebTokenStore) parse(value string) (*jwt.Token, error) {
	// parse the token and verify the signature
	token, err := jwt.ParseWithClaims(value, &ExtendedClaims{}, func(token *jwt.Token) (any, error) {
		return s.verifyKey, nil
	})
	if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, ErrInvalidTokenSignature
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *JsonWebTokenStore) Issue(_ context.Context, token *Token, ttl time.Duration) (string, error) {
	token.id = uuid.New().String()            // generate a new UUID for the token id
	token.issuedAt = time.Now().UTC()         // set the token create time
	token.expiresAt = token.issuedAt.Add(ttl) // set the token expiry time
	var claims = ExtendedClaims{              // create the token claims
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token.id,
			Issuer:    s.tokenIssuer,
			Audience:  jwt.ClaimStrings{s.tokenAudience},
			Subject:   token.Subject(),
			IssuedAt:  jwt.NewNumericDate(token.issuedAt),
			NotBefore: jwt.NewNumericDate(token.issuedAt),
			ExpiresAt: jwt.NewNumericDate(token.expiresAt),
		},
		Type:   token.ttype,
		Realm:  token.Realm(),
		Client: token.Client(),
		Scope:  strings.Join(token.Scope(), " "), // set the scope as a space-separated string
	}
	// sign the claims and return the signed token
	return jwt.NewWithClaims(s.signingMethod, claims).SignedString(s.signingKey)
}

func (s *JsonWebTokenStore) Renew(ctx context.Context, value string, ttl time.Duration) (string, error) {
	// parse the token and verify the signature
	token, err := s.parse(value)
	if err != nil {
		return "", err
	}
	// casting the token claims to the extended claims type
	// panic if the claims are not of the extended claims type
	claims := token.Claims.(*ExtendedClaims)
	if !strings.EqualFold(claims.Type, TokenTypeRefresh) {
		return "", ErrInvalidTokenType
	}
	return s.Issue(ctx, NewToken(TokenTypeBearer, claims.Realm, claims.Client, claims.Subject, strings.Split(claims.Scope, " ")), ttl)
}

func (s *JsonWebTokenStore) Verify(_ context.Context, value string) (*Token, error) {
	// parse the token and verify the signature
	token, err := s.parse(value)
	if err != nil {
		return nil, err
	}
	// casting the token claims to the extended claims type
	// panic if the claims are not of the extended claims type
	claims := token.Claims.(*ExtendedClaims)
	if !strings.EqualFold(claims.Type, TokenTypeBearer) {
		return nil, ErrInvalidTokenType
	}
	return &Token{
		id:        claims.ID,
		subject:   claims.Subject,
		issuedAt:  claims.IssuedAt.Time,
		expiresAt: claims.ExpiresAt.Time,
		client:    claims.Client,
		realm:     claims.Realm,
		scope:     strings.Split(claims.Scope, " "), // split the scope string into a slice of strings
	}, nil
}

func (s *JsonWebTokenStore) Revoke(_ context.Context, _ string) (*Token, error) {
	// Revoke is not supported for the JSON Web Token (JWT) token store.
	return nil, ErrUnsupportedOperation
}
