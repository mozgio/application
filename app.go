package app

import (
	"context"
	"io/fs"

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
	WithMiddleware(mw grpc.UnaryServerInterceptor) App
	WithGRPC(host string, port int) App
	WithHTTP(host string, port int) App
	WithDatabase(dsn string) App
	WithMigrations(migrationsFs fs.FS, pattern string) App
	WithService(factory ServiceFunc) App
	WithSwagger(contents []byte) App
	Listen()
}

type GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
type ServiceFunc func(Context) (*grpc.ServiceDesc, GatewayFunc, any)

type app struct {
	ctx                Context
	services           []ServiceFunc
	grpcListenerConfig listenerConfig
	httpListenerConfig listenerConfig
	withDatabase       bool
	databaseConfig     databaseConfig
	withMigrations     bool
	migrationsConfig   migrationsConfig
	middlewares        []grpc.UnaryServerInterceptor
	withSwagger        bool
	swaggerConfig      swaggerConfig
}

func (a *app) WithGRPC(host string, port int) App {
	a.grpcListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app) WithHTTP(host string, port int) App {
	a.httpListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app) WithService(factory ServiceFunc) App {
	a.services = append(a.services, factory)
	return a
}
