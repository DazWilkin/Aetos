package xxx

import (
	"crypto/md5"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"strings"
	"sync"

	pb "github.com/DazWilkin/Aetos/protos"

	"github.com/prometheus/client_golang/prometheus"
)

// TODO(dazwilkin) Should this be part of the Config struct?
var mu sync.RWMutex

// Config is a type that represents the current configuration of metrics|labels
type Config struct {
	Cardinality uint8
	NumLabels   uint8
	NumMetrics  uint8

	labels  []string
	metrics []string
}

// NewConfig is a function that creates a new Config
func NewConfig(cardinality, numLabels, numMetrics uint8) *Config {
	return &Config{
		Cardinality: cardinality,
		NumLabels:   numLabels,
		NumMetrics:  numMetrics,
		labels:      []string{},
		metrics:     []string{},
	}
}

// Get is a method that gets the Config as a protobuf message
func (c *Config) Get() *pb.AetosPublishRequest {
	return &pb.AetosPublishRequest{
		Labels:  c.labels,
		Metrics: c.metrics,
	}
}

// Set is a method that sets the Config using a protobuf message
func (c *Config) Set(rqst *pb.AetosPublishRequest) error {
	if len(rqst.Labels) == 0 {
		return fmt.Errorf("labels must be non-empty")
	}
	if len(rqst.Labels) > int(c.NumLabels) {
		return fmt.Errorf("too many labels (max: %d)", c.NumLabels)
	}
	if len(rqst.Metrics) == 0 {
		return fmt.Errorf("metrics must be non-empty")
	}
	if len(rqst.Metrics) > int(c.NumMetrics) {
		return fmt.Errorf("too many metrics (max: %d)", c.NumMetrics)
	}

	slog.Info("Set",
		"labels", strings.Join(rqst.Labels, ","),
		"metrics", strings.Join(rqst.Metrics, ","),
	)

	mu.Lock()
	c.labels = rqst.Labels
	c.metrics = rqst.Metrics
	mu.Unlock()

	return nil
}

// Descs is a method that represents the Config as a slice of prometheus.Desc
// This method is used by the Metrics methods and by the collector's Describe method
func (c *Config) Descs() []*prometheus.Desc {
	mu.RLock()
	result := make([]*prometheus.Desc, len(c.metrics))

	for i, metric := range c.metrics {
		result[i] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, metric),
			fmt.Sprintf("a randomly generated metric with %d labels with cardinality %d", len(c.labels), c.Cardinality),
			c.labels,
			nil,
		)
	}

	mu.RUnlock()
	return result
}

// Metrics is a method that represents the Config as a slice of prometheus.Metric
// This method is used by the collector's Collect method
func (c *Config) Metrics() []prometheus.Metric {
	mu.RLock()
	result := make([]prometheus.Metric, len(c.metrics)*int(c.Cardinality))

	// Enumerate each of the metrics (represented by prometheus.Desc)
	for i, desc := range c.Descs() {
		// For each cardinality
		// Effectively multiple the number of metrics
		for j := uint8(0); j < c.Cardinality; j++ {
			measurement := rand.Float64()

			// Create a new set of label values for this measurement
			labelValues := make([]string, len(c.labels))
			for k := 0; k < len(c.labels); k++ {
				// Value should be predicatable|repeatable
				// Want the same label values across different metrics
				// hash some combination of the label name and the cardinality
				value := fmt.Sprintf("%s-%d", c.labels[k], j)
				labelValues[k] = value // hash(value)
			}

			result[i*int(c.Cardinality)+int(j)] = prometheus.MustNewConstMetric(
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

// hash is a function that generates the MD5 hash of a string
func hash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", (h.Sum(nil)))
}
