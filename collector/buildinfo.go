package collector

import (
	"github.com/DazWilkin/Aetos/opts"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*BuildInfoCollector)(nil)

// BuildInfoCollector collects metrics, mostly runtime, about this exporter in general.
type BuildInfoCollector struct {
	opts      opts.Build
	startTime *prometheus.Desc
	buildInfo *prometheus.Desc
}

// NewBuildInfoCollector returns a new ExporterCollector.
func NewBuildInfoCollector(opts *opts.Build) *BuildInfoCollector {
	return &BuildInfoCollector{

		startTime: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "start_time"),
			"Exporter start time in Unix epoch seconds",
			nil,
			nil,
		),
		buildInfo: prometheus.NewDesc(
			prometheus.BuildFQName(opts.Namespace, opts.Subsystem, "build_info"),
			"A metric with a constant '1' value labeled by OS version, Go version, and the Git commit of the exporter",
			[]string{"os_version", "go_version", "git_commit"},
			nil,
		),
	}
}

// Collect implements Prometheus' Collector interface and is used to collect metrics
func (c *BuildInfoCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		c.startTime,
		prometheus.GaugeValue,
		float64(c.opts.StartTime),
	)
	ch <- prometheus.MustNewConstMetric(
		c.buildInfo,
		prometheus.CounterValue,
		1.0,
		c.opts.OSVersion, c.opts.GoVersion, c.opts.GitCommit,
	)
}

// Describe implements Prometheus' Collector interface and is used to describe metrics
func (c *BuildInfoCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.startTime
}
