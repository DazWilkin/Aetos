package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/DazWilkin/Aetos/collector"
	"github.com/DazWilkin/Aetos/handler"
	"github.com/DazWilkin/Aetos/opts"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace string = "aetos"
	subsystem string = "collector"
)
const (
	maxCardinality uint = 10
	maxLabels      uint = 5
	maxMetrics     uint = 250
)

var (
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	cardinality = flag.Uint("cardinality", 3, "Number of label values")
	endpoint    = flag.String("endpoint", ":8080", "Endpoint of ")
	labels      = flag.Uint("labels", 2, "Number of Labels")
	metrics     = flag.Uint("metrics", 5, "Number of Metrics")
	path        = flag.String("path", "metrics", "Path on which metrics will be served")
)

func main() {
	slog.Info("Configuration",
		"GitCommit", GitCommit,
		"GoVersion", GoVersion,
		"OSVersion", OSVersion,
		"StartTime", StartTime,
	)

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

	opts := opts.NewOpts(namespace, subsystem)
	optsBuild := opts.NewBuildOpts(GitCommit, GoVersion, OSVersion, StartTime)
	optsAetos := opts.NewAetosOpts(uint8(*cardinality), uint8(*labels), uint8(*metrics))

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewBuildInfoCollector(optsBuild))
	registry.MustRegister(collector.NewAetosCollector(optsAetos))

	mux := http.NewServeMux()

	// Index handler
	mux.HandleFunc("/", index)

	// Publish handler
	mux.Handle("/publish", handler.NewPublisher(optsAetos))

	// z-pages
	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/statusz", http.HandlerFunc(statusz))
	mux.Handle("/varz", handler.NewVarz(optsAetos))

	promOpts := promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}

	mux.Handle(fmt.Sprintf("/%s", *path), promhttp.HandlerFor(registry, promOpts))

	http.ListenAndServe(*endpoint, mux)
}
