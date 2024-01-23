package xxx

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var mu sync.RWMutex

type Foo struct {
	Cardinality uint8
	NumLabels   uint8
	NumMetrics  uint8

	labels  []string
	metrics []string
}

func NewFoo(cardinality, numLabels, numMetrics uint8) *Foo {
	return &Foo{
		Cardinality: cardinality,
		NumLabels:   numLabels,
		NumMetrics:  numMetrics,
		labels:      []string{},
		metrics:     []string{},
	}
}
func (f *Foo) Update(labels, metrics []string) error {
	if len(labels) == 0 {
		return fmt.Errorf("labels must be non-empty")
	}
	if len(labels) > int(f.NumLabels) {
		return fmt.Errorf("too many labels (max: %d)", f.NumLabels)
	}
	if len(metrics) == 0 {
		return fmt.Errorf("metrics must be non-empty")
	}
	if len(metrics) > int(f.NumMetrics) {
		return fmt.Errorf("too many metrics (max: %d)", f.NumMetrics)
	}

	mu.Lock()
	f.labels = labels
	f.metrics = metrics
	mu.Unlock()

	return nil
}
func (f *Foo) Descs() []*prometheus.Desc {
	result := make([]*prometheus.Desc, len(f.metrics))

	for i, metric := range f.metrics {
		result[i] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metric),
			fmt.Sprintf("a randomly generated metric with %d labels with cardinality %d", len(f.labels), f.Cardinality),
			f.labels,
			nil,
		)
	}

	return result
}
func (f *Foo) Metrics() []prometheus.Metric {
	result := make([]prometheus.Metric, len(f.metrics)*int(f.Cardinality))

	for i, desc := range f.Descs() {
		for j := uint8(0); j < f.Cardinality; j++ {
			value := rand.Float64()

			labelValues := make([]string, len(f.labels))
			for k := 0; k < len(f.labels); k++ {
				labelValues[k] = hash(
					fmt.Sprintf(
						"%s-%d",
						f.labels[k],
						j,
					),
				)
			}

			result[i*int(f.Cardinality)+int(j)] = prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				value,
				labelValues...,
			)
		}
	}

	return result
}

func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", (h.Sum(nil)))
}
