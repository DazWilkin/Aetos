package collector

import (
	"github.com/DazWilkin/Aetos/xxx"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that AetosCollector implements Prometheus' collector interface
var _ prometheus.Collector = (*AetosCollector)(nil)

// AetosCollector represents Aetos' dynamic metrics|labels
type AetosCollector struct {
	Config *xxx.Config
}

// NewAetosCollector is a function that creates a new AetosCollector
func NewAetosCollector(foo *xxx.Config) *AetosCollector {
	return &AetosCollector{
		Config: foo,
	}
}

// Collect is a method that implements Prometheus' Collector interface and is used to collect metrics
func (c *AetosCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.Config.Metrics() {
		ch <- metric
	}
}

// Describe is a method that implements Prometheus' Collector interface and is used to describe metrics
func (c *AetosCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.Config.Descs() {
		ch <- desc
	}
}
