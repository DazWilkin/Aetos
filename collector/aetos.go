package collector

import (
	"github.com/DazWilkin/Aetos/opts"

	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that AetosCollector implements Prometheus' collector interface
var _ prometheus.Collector = (*AetosCollector)(nil)

// AetosCollector represents Aetos' dynamic metrics|labels
type AetosCollector struct {
	opts *opts.Aetos
}

// NewAetosCollector is a function that creates a new AetosCollector
func NewAetosCollector(opts *opts.Aetos) *AetosCollector {
	return &AetosCollector{
		opts: opts,
	}
}

// Collect is a method that implements Prometheus' Collector interface and is used to collect metrics
func (c *AetosCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.opts.Metrics() {
		ch <- metric
	}
}

// Describe is a method that implements Prometheus' Collector interface and is used to describe metrics
func (c *AetosCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.opts.Descs() {
		ch <- desc
	}
}
