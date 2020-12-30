package main

import (
	"net/http"
	"os"
	"time"

	"github.com/auifzysr/envii_exporter/pkg/exporter"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
)

var metricsPath = "/metrics"

func main() {
	// TOOD: parameterize
	address := "http://192.168.11.12:9000"
	timeout := time.Duration(5)
	logger := promlog.New(&promlog.Config{})
	listenAddress := "0.0.0.0:8001"

	level.Info(logger).Log("msg", "Starting memcached_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	prometheus.MustRegister(version.NewCollector("envii_exporter"))
	prometheus.MustRegister(exporter.New(address, timeout, logger))

	http.Handle(metricsPath, promhttp.Handler())

	level.Info(logger).Log("msg", "Listening on address", "address", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", "Error running http server", "err", err)
		os.Exit(1)
	}
}
