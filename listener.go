package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mozgio/app/swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type listenerConfig struct {
	host string
	port int
}

func (lc listenerConfig) address() string {
	return fmt.Sprintf("%s:%d", lc.host, lc.port)
}
func (a *app) WithMiddleware(mw grpc.UnaryServerInterceptor) App {
	a.middlewares = append(a.middlewares, mw)
	return a
}

func (a *app) Listen() {
	if a.withDatabase {
		a.ctx = connectToDatabase(a.ctx, a.databaseConfig)
	}
	if a.withMigrations {
		migrateDatabase(a.ctx, a.migrationsConfig)
	}
	grpcServerAddr := a.grpcListenerConfig.address()
	httpServerAddr := a.httpListenerConfig.address()

	// Configure GRPC-GATEWAY
	gatewayMux := runtime.NewServeMux(
		runtime.WithOutgoingHeaderMatcher(outgoingHeaderMatcher),
	)
	gatewayOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Configure GRPC server

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(a.middlewares...),
	)
	for _, factory := range a.services {
		desc, gw, impl := factory(a.ctx)

		grpcServer.RegisterService(desc, impl)
		if err := gw(a.ctx, gatewayMux, a.grpcListenerConfig.address(), gatewayOpts); err != nil {
			a.ctx.Log().Fatal("failed to setup grpc-gateway",
				zap.Error(err))
		}
	}

	// GRPC service listener
	go func() {
		a.ctx.Log().Info("starting grpc service",
			zap.String("addr", a.grpcListenerConfig.address()))
		grpcListener, err := net.Listen("tcp", grpcServerAddr)
		if err != nil {
			a.ctx.Log().Fatal("failed to listen grpc",
				zap.Int("port", a.grpcListenerConfig.port),
				zap.Error(err))
		}
		if err = grpcServer.Serve(grpcListener); err != nil {
			a.ctx.Log().Fatal("failed to serve grpc",
				zap.Error(err))
		}
	}()

	// GRPC-GATEWAY service listener
	var handler http.Handler

	if a.withSwagger {
		handler = swagger.WithSwagger(a.swaggerConfig.fileContext, gatewayMux)
	} else {
		handler = gatewayMux
	}

	gatewayServer := &http.Server{
		Addr:    httpServerAddr,
		Handler: handler,
	}
	go func() {
		a.ctx.Log().Info("starting http service",
			zap.String("addr", gatewayServer.Addr))
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.ctx.Log().Fatal("failed to listen http",
				zap.Int("port", a.httpListenerConfig.port),
				zap.Error(err))
		}
	}()

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

		grpcServer.GracefulStop()

		if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
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

//func (a *app) withCustomPaths(mux http.Handler) http.Handler {
//	staticServer := http.FileServer(http.FS(a.swaggerFile.FS))
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if strings.HasPrefix(r.URL.Path, "/swagger") {
//			staticServer.ServeHTTP(w, r)
//			return
//		}
//		mux.ServeHTTP(w, r)
//	})
//}

func outgoingHeaderMatcher(key string) (string, bool) {
	fmt.Println("out", key)
	if key == "set-cookie" {
		return key, true
	}
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}
