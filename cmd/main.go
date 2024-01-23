package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/DazWilkin/Aetos/collector"
	"github.com/DazWilkin/Aetos/handler"
	"github.com/DazWilkin/Aetos/xxx"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	maxCardinality uint = 10
	maxLabels      uint = 5
	maxMetrics     uint = 250
)

var (
	// Path will be set to the GitHub repo
	// Version
	Version string
	// Checksum is the git commit value and is expected to be set during build
	Checksum string
)
var (
	cardinality = flag.Uint("cardinality", 3, "Number of label values")
	endpoint    = flag.String("endpoint", ":8080", "Endpoint of ")
	labels      = flag.Uint("labels", 2, "Number of Labels")
	metrics     = flag.Uint("metrics", 5, "Number of Metrics")
	path        = flag.String("path", "metrics", "Path on which metrics will be served")
)

func main() {
	flag.Parse()

	if *cardinality > maxCardinality {
		slog.Info("Flag `--cardinality` too large", "max", maxCardinality)
		os.Exit(1)
	}
	if *labels > maxLabels {
		slog.Info("Flag `--max_labels` too large", "max", maxLabels)
		os.Exit(1)
	}
	if *metrics > maxMetrics {
		slog.Info("Flag `--max_metrics` too large", "max", maxMetrics)
		os.Exit(1)
	}
	if *path == "" {
		slog.Info("Flag `--path` must be non-empty")
		os.Exit(1)
	}

	slog.Info("Configuration",
		"cardinality", *cardinality,
		"labels", *labels,
		"metrics", *metrics,
		"path", *path,
	)

	config := xxx.NewConfig(uint8(*cardinality), uint8(*labels), uint8(*metrics))

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		// collectors.NewGoCollector(
		// 	collectors.WithGoCollectorRuntimeMetrics(),
		// ),
		collectors.NewBuildInfoCollector(),
		// collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	registry.MustRegister(collector.NewAetosCollector(config))

	mux := http.NewServeMux()
	mux.Handle("/publish", handler.NewPublisher(config))

	// z-pages
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/statusz", http.HandlerFunc(statusz))
	mux.Handle("/varz", handler.NewVarz(config))

	opts := promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}

	mux.Handle(fmt.Sprintf("/%s", *path), promhttp.HandlerFor(registry, opts))

	http.ListenAndServe(*endpoint, mux)
}
