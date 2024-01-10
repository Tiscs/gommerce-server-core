package events

import (
	"github.com/choral-io/gommerce-server-core/config"
	"github.com/nats-io/nats.go"
)

func NewNATSConn(cfg config.ServerNATSConfig) (*nats.Conn, error) {
	return nats.Connect(cfg.GetURL(), nats.NoEcho())
}
