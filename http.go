package Application

import (
	"errors"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/mozgio/application/Swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func (a *app[TConfig, Database]) WithHTTP(host string, port int) App[TConfig, Database] {
	a.withHTTP = true
	a.httpListenerConfig = listenerConfig{host, port}
	return a
}

func (a *app[TConfig, Database]) serveHTTP() {
	var handler http.Handler = a.gatewayMux

	if a.withSwagger {
		handler = Swagger.WithSwagger(a.swaggerConfig.fileContext, handler)
	}

	if a.withMetrics {
		handler = wrapWithMetrics(a.serverMetrics, handler)
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

func wrapWithMetrics(serverMetrics *prometheus.ServerMetrics, nextHandler http.Handler) http.Handler {
	promHandler := promhttp.Handler()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/metrics") {
			promHandler.ServeHTTP(w, r)
			return
		}
		nextHandler.ServeHTTP(w, r)
	})
}
