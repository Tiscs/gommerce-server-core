package events

import (
	"github.com/choral-io/gommerce-server-core/config"
	"github.com/nats-io/nats.go"
)

func NewNATSConn(cfg config.ServerNATSConfig) (*nats.Conn, error) {
	opts := []nats.Option{}
	if cfg.GetNoEcho() {
		opts = append(opts, nats.NoEcho())
	}
	return nats.Connect(cfg.GetSeedURL(), opts...)
}
