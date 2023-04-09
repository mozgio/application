package Metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metric interface {
	Collector() prometheus.Collector
}

func Counter(name, help string, labelKVPairs ...string) prometheus.Collector {
	if len(labelKVPairs)%2 != 0 {
		panic("[Metrics.Counter] invalid amount of key-value pairs")
	}
	labels := make(map[string]string, len(labelKVPairs))
	for i := 0; i < len(labelKVPairs); i = i + 2 {
		labels[labelKVPairs[i]] = labelKVPairs[i+1]
	}
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   "mozg",
		Subsystem:   "application",
		Name:        name,
		Help:        help,
		ConstLabels: labels,
	})
}
