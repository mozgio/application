# application

Application package.

## example

```go
package main

import (
	"time"

	Application "github.com/mozgio/application"
	"go.uber.org/zap"
)

func main() {
	cfg := Application.MustParseConfig[Config]()

	run := func(ctx Application.Context[Config, any]) {
		timer := time.NewTimer(time.Second * 5)
		for t := range timer.C {
			ctx.Log().Info("time", zap.Time("time", t))
		}
	}

	Application.
		New[Config, any]().
		WithConfig(cfg).
		WithGRPC(cfg.GRPCAddress()).
		WithHTTP(cfg.HTTPAddress()).
		WithRunner(run).
		Listen()
}

// Config struct.
type Config struct {
	GRPCHost    string `env:"GRPC_HOST" envDefault:"0.0.0.0"`
	GRPCPort    int    `env:"GRPC_PORT" envDefault:"9000"`
	HTTPHost    string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	HTTPPort    int    `env:"HTTP_PORT" envDefault:"9080"`
	DatabaseDsn string `env:"DATABASE_DSN" envDefault:"root:password@tcp(localhost:3306)/mozg"`
	NatsURL     string `env:"NATS_URL" envDefault:"demo.nats.io:4222"`
}

func (c *Config) GRPCAddress() (string, int) {
	return c.GRPCHost, c.GRPCPort
}

func (c *Config) HTTPAddress() (string, int) {
	return c.HTTPHost, c.HTTPPort
}

```