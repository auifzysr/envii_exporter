package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/auifzysr/envii_exporter/pkg/exporter"
	"github.com/go-kit/kit/log/level"
	"github.com/jamiealquiza/envy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
)

var metricsPath = "/metrics"

func main() {
	var nodeEndpoint = flag.String("node_endpoint", "localhost:9000", "The endpoint to which the exporter access to get metrics.")
	// TODO: enable timeout
	var timeout = flag.Int("timeout", 15, "Seconds the exporter waits to get response.")
	var listenAddress = flag.String("listen_address", "0.0.0.0:8001", "The address on which the exporter listens for connections.")
	// the alternative way to pass parameters with environment variables of which the name begins with "ENVII"
	envy.Parse("ENVII")
	flag.Parse()

	// initializes Go kit loggers across Prometheus components
	// https://godoc.org/github.com/prometheus/common/promlog
	var logger = promlog.New(&promlog.Config{})
	// TODO: embed build info
	level.Info(logger).Log("msg", "Starting envii_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	// registers Collectors
	prometheus.MustRegister(version.NewCollector("envii_exporter"))
	prometheus.MustRegister(exporter.New(*nodeEndpoint, time.Duration(*timeout), logger))

	// registers at which to expose metrics
	http.Handle(metricsPath, promhttp.Handler())

	// listens connections from Prometheus
	level.Info(logger).Log("msg", "Listening on address", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error running http server", "err", err)
		os.Exit(1)
	}
}
