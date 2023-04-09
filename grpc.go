package Application

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func (a *app[TConfig, TDatabase]) WithGRPC(host string, port int) App[TConfig, TDatabase] {
	a.withGRPC = true
	a.grpcListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app[TConfig, TDatabase]) WithGRPCServerOption(opts ...grpc.ServerOption) App[TConfig, TDatabase] {
	a.grpcServerOpts = append(a.grpcServerOpts, opts...)
	return a
}

func (a *app[TConfig, TDatabase]) WithUnaryInterceptor(interceptor grpc.UnaryServerInterceptor) App[TConfig, TDatabase] {
	a.unaryInterceptors = append(a.unaryInterceptors, interceptor)
	return a
}

func (a *app[TConfig, TDatabase]) WithStreamInterceptor(interceptor grpc.StreamServerInterceptor) App[TConfig, TDatabase] {
	a.streamInterceptors = append(a.streamInterceptors, interceptor)
	return a
}

func (a *app[TConfig, TDatabase]) serveGRPC() {
	a.grpcServer = grpc.NewServer(append(a.grpcServerOpts,
		grpc.ChainUnaryInterceptor(a.unaryInterceptors...),
		grpc.ChainStreamInterceptor(a.streamInterceptors...),
	)...)

	if a.withMetrics {
		a.serverMetrics.InitializeMetrics(a.grpcServer)
	}

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
