package Application

import (
	"github.com/prometheus/client_golang/prometheus"
)

func (a *app[TConfig, TDatabase]) WithMetrics(metrics ...prometheus.Collector) App[TConfig, TDatabase] {
	a.withMetrics = true
	a.metrics = append(a.metrics, metrics...)
	return a
}
