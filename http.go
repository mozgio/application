package app

import (
	"errors"
	"net/http"

	"github.com/mozgio/app/swagger"
	"go.uber.org/zap"
)

func (a *app) WithHTTP(host string, port int) App {
	a.withHTTP = true
	a.httpListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app) serveHTTP() {
	var handler http.Handler = a.gatewayMux

	if a.withSwagger {
		handler = swagger.WithSwagger(a.swaggerConfig.fileContext, handler)
	}

	a.httpServer = &http.Server{
		Addr:    a.httpListenerConfig.address(),
		Handler: handler,
	}

	a.ctx.Log().Info("starting http service",
		zap.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.ctx.Log().Fatal("failed to listen http",
			zap.Int("port", a.httpListenerConfig.port),
			zap.Error(err))
	}
}
