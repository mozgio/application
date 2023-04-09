package Application

import (
	"context"
	"io/fs"
	"net/http"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mozgio/application/Metrics"
	Database "github.com/mozgio/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func New[TConfig ConfigType, TDatabase DatabaseType]() App[TConfig, TDatabase] {
	ctx := newContext[TConfig, TDatabase]()

	serverMetrics := grpcPrometheus.NewServerMetrics()

	a := &app[TConfig, TDatabase]{
		ctx:           ctx,
		serverMetrics: serverMetrics,
		unaryInterceptors: []grpc.UnaryServerInterceptor{
			serverMetrics.UnaryServerInterceptor(),
		},
		streamInterceptors: []grpc.StreamServerInterceptor{
			serverMetrics.StreamServerInterceptor(),
		},
		gatewayGRPCOpts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}

	return a
}

type App[TConfig ConfigType, TDatabase DatabaseType] interface {
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
	WithMetrics(...Metrics.Metric) App[TConfig, TDatabase]
	Listen()
}

type ConfigType any
type DatabaseType any

type GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
type ServiceFunc[TConfig ConfigType, TDatabase DatabaseType] func(Context[TConfig, TDatabase]) (*grpc.ServiceDesc, GatewayFunc, any)

type app[TConfig ConfigType, TDatabase DatabaseType] struct {
	ctx *appContext[TConfig, TDatabase]

	withMetrics   bool
	metrics       []Metrics.Metric
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
}
