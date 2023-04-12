package application

import (
	"context"
	"io/fs"
	"net/http"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	Database "github.com/mozgio/database"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New[TConfig any, TDatabase any]() App[TConfig, TDatabase] {
	ctx := newContext[TConfig, TDatabase]()

	serverMetrics := grpcPrometheus.NewServerMetrics()
	prometheus.MustRegister(serverMetrics)
	logInterceptor := loggerInterceptor(ctx.logger)

	a := &app[TConfig, TDatabase]{
		ctx:           ctx,
		serverMetrics: serverMetrics,
		unaryInterceptors: []grpc.UnaryServerInterceptor{
			logging.UnaryServerInterceptor(logInterceptor),
			serverMetrics.UnaryServerInterceptor(),
		},
		streamInterceptors: []grpc.StreamServerInterceptor{
			logging.StreamServerInterceptor(logInterceptor),
			serverMetrics.StreamServerInterceptor(),
		},
		gatewayGRPCOpts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}

	return a
}

type App[TConfig any, TDatabase any] interface {
	WithConfig(cfg *TConfig) App[TConfig, TDatabase]
	WithGRPC(host string, port int) App[TConfig, TDatabase]
	WithHTTP(host string, port int) App[TConfig, TDatabase]
	WithGRPCServerOption(opts ...grpc.ServerOption) App[TConfig, TDatabase]
	WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) App[TConfig, TDatabase]
	WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) App[TConfig, TDatabase]
	WithDatabase(driver Database.Driver[TDatabase]) App[TConfig, TDatabase]
	WithMigrations(migrationsFs fs.FS, pattern string) App[TConfig, TDatabase]
	WithService(factory ServiceFunc[TConfig, TDatabase]) App[TConfig, TDatabase]
	WithSwagger(contents []byte) App[TConfig, TDatabase]
	WithMetrics(...prometheus.Collector) App[TConfig, TDatabase]
	WithNats(uri string, opts ...nats.Option) App[TConfig, TDatabase]
	WithRunner(runnerFunc RunnerFunc[TConfig, TDatabase]) App[TConfig, TDatabase]
	Listen()
}

//type any any
//type any any

type GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
type ServiceFunc[TConfig any, TDatabase any] func(Context[TConfig, TDatabase]) (*grpc.ServiceDesc, GatewayFunc, any)
type RunnerFunc[TConfig any, TDatabase any] func(Context[TConfig, TDatabase])

type app[TConfig any, TDatabase any] struct {
	ctx *appContext[TConfig, TDatabase]

	withMetrics   bool
	metrics       []prometheus.Collector
	serverMetrics *grpcPrometheus.ServerMetrics

	gatewayMux      *runtime.ServeMux
	gatewayGRPCOpts []grpc.DialOption
	gatewayMuxOpts  []runtime.ServeMuxOption

	grpcServer *grpc.Server
	httpServer *http.Server

	services []ServiceFunc[TConfig, TDatabase]

	grpcServerOpts     []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	withGRPC           bool
	grpcListenerConfig listenerConfig

	withHTTP           bool
	httpListenerConfig listenerConfig

	withDatabase   bool
	databaseDriver Database.Driver[TDatabase]

	withMigrations   bool
	migrationsConfig migrationsConfig

	withSwagger   bool
	swaggerConfig swaggerConfig

	withNats   bool
	natsConfig natsConfig

	withRunners bool
	runners     []RunnerFunc[TConfig, TDatabase]
}
