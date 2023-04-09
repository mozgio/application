package Application

import (
	"github.com/nats-io/nats.go"
)

func (a *app[TConfig, TDatabase]) WithNats(uri string, opts ...nats.Option) App[TConfig, TDatabase] {
	a.withNats = true
	a.natsConfig = natsConfig{
		uri: uri,
	}
	return a
}

type natsConfig struct {
	uri  string
	opts []nats.Option
}
