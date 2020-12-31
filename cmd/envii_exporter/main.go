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
	var poll_address = flag.String("poll_address", "localhost:9000", "The address to which the exporter access to get metrics.")
	var timeout = flag.Int("timeout", 15, "Seconds the exporter waits to get response.")
	var listen_address = flag.String("listen_address", "0.0.0.0:8001", "The address on which the exporter listens for connections.")
	envy.Parse("ENVII")
	flag.Parse()

	var logger = promlog.New(&promlog.Config{})
	level.Info(logger).Log("msg", "Starting envii_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	prometheus.MustRegister(version.NewCollector("envii_exporter"))
	prometheus.MustRegister(exporter.New(*poll_address, time.Duration(*timeout), logger))

	http.Handle(metricsPath, promhttp.Handler())

	level.Info(logger).Log("msg", "Listening on address", "address", *listen_address)
	if err := http.ListenAndServe(*listen_address, nil); err != nil {
		level.Error(logger).Log("msg", "Error running http server", "err", err)
		os.Exit(1)
	}
}
