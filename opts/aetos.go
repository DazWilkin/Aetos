package opts

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"sync"

	"github.com/DazWilkin/Aetos/api/v1alpha1"

	"github.com/prometheus/client_golang/prometheus"
)

// Aetos is a type that represents the current configuration of metrics|labels
type Aetos struct {
	mu sync.RWMutex

	namespace string
	subsystem string

	cardinality uint8
	NumLabels   uint8
	NumMetrics  uint8

	labels  []string
	metrics []string
}

// Get is a method that gets the Config as a protobuf message
func (a *Aetos) Get() *v1alpha1.AetosPublishRequest {
	return &v1alpha1.AetosPublishRequest{
		Labels:  a.labels,
		Metrics: a.metrics,
	}
}

// Set is a method that sets the Config using a protobuf message
func (a *Aetos) Set(rqst *v1alpha1.AetosPublishRequest) error {
	if len(rqst.Labels) == 0 {
		return fmt.Errorf("labels must be non-empty")
	}
	if len(rqst.Labels) > int(a.NumLabels) {
		return fmt.Errorf("too many labels (max: %d)", a.NumLabels)
	}
	if len(rqst.Metrics) == 0 {
		return fmt.Errorf("metrics must be non-empty")
	}
	if len(rqst.Metrics) > int(a.NumMetrics) {
		return fmt.Errorf("too many metrics (max: %d)", a.NumMetrics)
	}

	slog.Info("Set",
		"labels", strings.Join(rqst.Labels, ","),
		"metrics", strings.Join(rqst.Metrics, ","),
	)

	a.mu.Lock()
	a.labels = rqst.Labels
	a.metrics = rqst.Metrics
	a.mu.Unlock()

	return nil
}

// Descs is a method that represents the Config as a slice of prometheus.Desc
// This method is used by the Metrics methods and by the collector's Describe method
func (a *Aetos) Descs() []*prometheus.Desc {
	a.mu.RLock()
	result := make([]*prometheus.Desc, len(a.metrics))

	for i, metric := range a.metrics {
		result[i] = prometheus.NewDesc(
			prometheus.BuildFQName(a.namespace, a.subsystem, metric),
			fmt.Sprintf("a randomly generated metric with %d labels with cardinality %d", len(a.labels), a.cardinality),
			a.labels,
			nil,
		)
	}

	a.mu.RUnlock()
	return result
}

// Metrics is a method that represents the Config as a slice of prometheus.Metric
// This method is used by the collector's Collect method
func (a *Aetos) Metrics() []prometheus.Metric {
	a.mu.RLock()
	result := make([]prometheus.Metric, len(a.metrics)*int(a.cardinality))

	// Enumerate each of the metrics (represented by prometheus.Desc)
	for i, desc := range a.Descs() {
		// For each cardinality
		// Effectively multiple the number of metrics
		for j := uint8(0); j < a.cardinality; j++ {
			measurement := rand.Float64()

			// Create a new set of label values for this measurement
			labelValues := make([]string, len(a.labels))
			for k := 0; k < len(a.labels); k++ {
				// Value should be predicatable|repeatable
				// Want the same label values across different metrics
				// hash some combination of the label name and the cardinality
				value := fmt.Sprintf("%s-%d", a.labels[k], j)
				labelValues[k] = value // hash(value)
			}

			result[i*int(a.cardinality)+int(j)] = prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				measurement,
				labelValues...,
			)
		}
	}

	a.mu.RUnlock()
	return result
}
