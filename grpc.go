package app

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (a *app) WithGRPC(host string, port int) App {
	a.withGRPC = true
	a.grpcListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app) WithGRPCServerOption(opts ...grpc.ServerOption) App {
	a.grpcServerOpts = append(a.grpcServerOpts, opts...)
	return a
}

func (a *app) WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) App {
	a.unaryInterceptors = append(a.unaryInterceptors, interceptor)
	return a
}

func (a *app) WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) App {
	a.streamInterceptors = append(a.streamInterceptors, interceptor)
	return a
}

func (a *app) serveGRPC() {
	a.grpcServer = grpc.NewServer(append(a.grpcServerOpts,
		grpc.ChainUnaryInterceptor(a.unaryInterceptors...),
		grpc.ChainStreamInterceptor(a.streamInterceptors...),
	)...)

	a.configureServices()

	a.ctx.Log().Info("starting grpc service",
		zap.String("addr", a.grpcListenerConfig.address()))

	grpcListener, err := net.Listen("tcp", a.grpcListenerConfig.address())
	if err != nil {
		a.ctx.Log().Fatal("failed to listen grpc",
			zap.Int("port", a.grpcListenerConfig.port),
			zap.Error(err))
	}

	if err = a.grpcServer.Serve(grpcListener); err != nil {
		a.ctx.Log().Fatal("failed to serve grpc",
			zap.Error(err))
	}
}
