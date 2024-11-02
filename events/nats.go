package events

import (
	"github.com/nats-io/nats.go"

	"github.com/choral-io/gommerce-server-core/config"
)

func NewNATSConn(cfg config.ServerNATSConfig) (*nats.Conn, error) {
	var opts []nats.Option
	if cfg.GetNoEcho() {
		opts = append(opts, nats.NoEcho())
	}
	return nats.Connect(cfg.GetSeedURL(), opts...)
}
