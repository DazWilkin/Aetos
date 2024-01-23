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
	"github.com/DazWilkin/Aetos/xxx"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
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

	foo := xxx.NewFoo(uint8(*cardinality), uint8(*labels), uint8(*metrics))

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewAetosCollector(foo))

	mux := http.NewServeMux()
	mux.Handle("/publish", handler.NewPublisher(foo))
	mux.Handle("/healthz", http.HandlerFunc(healthz))

	opts := promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}

	mux.Handle(fmt.Sprintf("/%s", *path), promhttp.HandlerFor(registry, opts))

	http.ListenAndServe(*endpoint, mux)
}
