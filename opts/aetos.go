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

// TODO(dazwilkin) Should this be part of the Config struct?
var mu sync.RWMutex

// Aetos is a type that represents the current configuration of metrics|labels
type Aetos struct {
	namespace string
	subsystem string

	cardinality uint8
	NumLabels   uint8
	NumMetrics  uint8

	labels  []string
	metrics []string
}

// Get is a method that gets the Config as a protobuf message
func (o *Aetos) Get() *v1alpha1.AetosPublishRequest {
	return &v1alpha1.AetosPublishRequest{
		Labels:  o.labels,
		Metrics: o.metrics,
	}
}

// Set is a method that sets the Config using a protobuf message
func (o *Aetos) Set(rqst *v1alpha1.AetosPublishRequest) error {
	if len(rqst.Labels) == 0 {
		return fmt.Errorf("labels must be non-empty")
	}
	if len(rqst.Labels) > int(o.NumLabels) {
		return fmt.Errorf("too many labels (max: %d)", o.NumLabels)
	}
	if len(rqst.Metrics) == 0 {
		return fmt.Errorf("metrics must be non-empty")
	}
	if len(rqst.Metrics) > int(o.NumMetrics) {
		return fmt.Errorf("too many metrics (max: %d)", o.NumMetrics)
	}

	slog.Info("Set",
		"labels", strings.Join(rqst.Labels, ","),
		"metrics", strings.Join(rqst.Metrics, ","),
	)

	mu.Lock()
	o.labels = rqst.Labels
	o.metrics = rqst.Metrics
	mu.Unlock()

	return nil
}

// Descs is a method that represents the Config as a slice of prometheus.Desc
// This method is used by the Metrics methods and by the collector's Describe method
func (o *Aetos) Descs() []*prometheus.Desc {
	mu.RLock()
	result := make([]*prometheus.Desc, len(o.metrics))

	for i, metric := range o.metrics {
		result[i] = prometheus.NewDesc(
			prometheus.BuildFQName(o.namespace, o.subsystem, metric),
			fmt.Sprintf("a randomly generated metric with %d labels with cardinality %d", len(o.labels), o.cardinality),
			o.labels,
			nil,
		)
	}

	mu.RUnlock()
	return result
}

// Metrics is a method that represents the Config as a slice of prometheus.Metric
// This method is used by the collector's Collect method
func (o *Aetos) Metrics() []prometheus.Metric {
	mu.RLock()
	result := make([]prometheus.Metric, len(o.metrics)*int(o.cardinality))

	// Enumerate each of the metrics (represented by prometheus.Desc)
	for i, desc := range o.Descs() {
		// For each cardinality
		// Effectively multiple the number of metrics
		for j := uint8(0); j < o.cardinality; j++ {
			measurement := rand.Float64()

			// Create a new set of label values for this measurement
			labelValues := make([]string, len(o.labels))
			for k := 0; k < len(o.labels); k++ {
				// Value should be predicatable|repeatable
				// Want the same label values across different metrics
				// hash some combination of the label name and the cardinality
				value := fmt.Sprintf("%s-%d", o.labels[k], j)
				labelValues[k] = value // hash(value)
			}

			result[i*int(o.cardinality)+int(j)] = prometheus.MustNewConstMetric(
				desc,
				prometheus.GaugeValue,
				measurement,
				labelValues...,
			)
		}
	}

	mu.RUnlock()
	return result
}
