package secure

import (
	"encoding/json"
	"time"
)

type tokenData struct {
	TType     string    `json:"ttype,omitempty"`
	Realm     string    `json:"realm,omitempty"`
	Client    string    `json:"client,omitempty"`
	Subject   string    `json:"subject,omitempty"`
	Scope     []string  `json:"scope,omitempty"`
	IssuedAt  time.Time `json:"issued_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

func (t *Token) MarshalJSON() ([]byte, error) {
	data := &tokenData{
		TType:     t.ttype,
		Realm:     t.realm,
		Client:    t.client,
		Subject:   t.subject,
		Scope:     t.scope,
		IssuedAt:  t.issuedAt,
		ExpiresAt: t.expiresAt,
	}
	return json.Marshal(data)
}

func (t *Token) UnmarshalJSON(data []byte) error {
	var td tokenData
	if err := json.Unmarshal(data, &td); err != nil {
		return err
	}
	t.ttype = td.TType
	t.realm = td.Realm
	t.client = td.Client
	t.subject = td.Subject
	t.scope = td.Scope
	t.issuedAt = td.IssuedAt
	t.expiresAt = td.ExpiresAt
	return nil
}
