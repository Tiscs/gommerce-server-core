package secure

// Identity represents the result of authentication.
type Identity struct {
	schema string
	token  *Token
}

// NewIdentity returns a new Identity with the given schema and token.
func NewIdentity(schema string, token *Token) *Identity {
	return &Identity{
		schema: schema,
		token:  token,
	}
}

// Schema returns the schema of the identity.
func (i *Identity) Schema() string {
	return i.schema
}

// Token returns the token of the identity.
func (i *Identity) Token() *Token {
	return i.token
}
