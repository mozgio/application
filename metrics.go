package Application

import (
	"github.com/mozgio/application/Metrics"
)

func (a *app[TConfig, TDatabase]) WithMetrics(metrics ...Metrics.Metric) App[TConfig, TDatabase] {
	a.withMetrics = true
	a.metrics = append(a.metrics, metrics...)
	return a
}
