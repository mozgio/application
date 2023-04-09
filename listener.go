package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
)

type listenerConfig struct {
	host string
	port int
}

func (lc listenerConfig) address() string {
	return fmt.Sprintf("%s:%d", lc.host, lc.port)
}

func (a *app) Listen() {
	if a.withDatabase {
		a.ctx = connectToDatabase(a.ctx, a.databaseConfig)
		if a.withMigrations {
			migrateDatabase(a.ctx, a.migrationsConfig)
		}
	}

	a.gatewayMux = runtime.NewServeMux(a.gatewayMuxOpts...)

	if a.withGRPC {
		go a.serveGRPC()
	}

	if a.withHTTP {
		go a.serveHTTP()
	}

	done := make(chan struct{}, 1)
	defer close(done)
	go func() {
		shutdownCtx, cancel := context.WithTimeout(context.TODO(), time.Second*30)
		defer cancel()

		shutdown := make(chan os.Signal, 1)
		defer close(shutdown)
		signal.Notify(shutdown, os.Interrupt, os.Kill)
		sig := <-shutdown

		a.ctx.Log().Info("graceful shutdown start",
			zap.String("signal", sig.String()))

		a.grpcServer.GracefulStop()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			a.ctx.Log().Error("failed to shutdown http gracefully",
				zap.Error(err))
		}

		var closers []func() error
		if a.withDatabase {
			closers = append(closers, a.ctx.DbConn().Close)
		}

		for _, closer := range closers {
			if err := closer(); err != nil {
				a.ctx.Log().Warn("failed to close on shutdown",
					zap.Error(err))
			}
		}

		done <- struct{}{}
	}()

	<-done
	a.ctx.Log().Info("graceful shutdown complete")
}
