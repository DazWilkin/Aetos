package collector

import (
	"github.com/DazWilkin/Aetos/xxx"
	"github.com/prometheus/client_golang/prometheus"
)

// Ensure that AetosCollector implements Prometheus' collector interface
var _ prometheus.Collector = (*AetosCollector)(nil)

type AetosCollector struct {
	Foo *xxx.Foo
}

func NewAetosCollector(foo *xxx.Foo) *AetosCollector {
	return &AetosCollector{
		Foo: foo,
	}
}

func (c *AetosCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range c.Foo.Metrics() {
		ch <- metric
	}
}
func (c *AetosCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range c.Foo.Descs() {
		ch <- desc
	}
}
