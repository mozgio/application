package app

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func New() App {
	ctx := newContext()

	a := &app{
		ctx: ctx,
	}

	return a
}

type App interface {
	WithGRPC(host string, port int) App
	WithHTTP(host string, port int) App
	WithGRPCServerOption(opts ...grpc.ServerOption) App
	WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) App
	WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) App
	WithDatabase(dsn string) App
	WithMigrations(migrationsFs fs.FS, pattern string) App
	WithService(factory ServiceFunc) App
	WithSwagger(contents []byte) App
	Listen()
}

type GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
type ServiceFunc func(Context) (*grpc.ServiceDesc, GatewayFunc, any)

type app struct {
	ctx Context

	gatewayMux      *runtime.ServeMux
	gatewayGRPCOpts []grpc.DialOption
	gatewayMuxOpts  []runtime.ServeMuxOption

	grpcServer *grpc.Server
	httpServer *http.Server

	services []ServiceFunc

	grpcServerOpts     []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor

	withGRPC           bool
	grpcListenerConfig listenerConfig

	withHTTP           bool
	httpListenerConfig listenerConfig

	withDatabase   bool
	databaseConfig databaseConfig

	withMigrations   bool
	migrationsConfig migrationsConfig

	withSwagger   bool
	swaggerConfig swaggerConfig
}
