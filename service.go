package app

import (
	"go.uber.org/zap"
)

func (a *app) WithService(factory ServiceFunc) App {
	a.services = append(a.services, factory)
	return a
}

func (a *app) configureServices() {
	for _, factory := range a.services {
		desc, gw, impl := factory(a.ctx)

		a.grpcServer.RegisterService(desc, impl)
		if err := gw(a.ctx, a.gatewayMux, a.grpcListenerConfig.address(), a.gatewayGRPCOpts); err != nil {
			a.ctx.Log().Fatal("failed to setup grpc-gateway",
				zap.Error(err))
		}
	}
}
