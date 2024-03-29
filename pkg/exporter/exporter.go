package exporter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	nodeEndpoint string
	timeout      time.Duration
	logger       log.Logger

	up      *prometheus.Desc
	version *prometheus.Desc

	temperature *prometheus.Desc
	humidity    *prometheus.Desc
	pressure    *prometheus.Desc
}

const (
	Namespace = "envii"
)

// builds an exporter object
func New(nodeEndpoint string, timeout time.Duration, logger log.Logger) *Exporter {
	return &Exporter{
		nodeEndpoint: nodeEndpoint,
		timeout:      timeout,
		logger:       logger,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "up"),
			"up",
			nil,
			nil,
		),
		version: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "version"),
			"version",
			nil,
			nil,
		),
		temperature: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "temperature"),
			"temperature",
			nil,
			nil,
		),
		humidity: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "humidity"),
			"humidity",
			nil,
			nil,
		),
		pressure: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "pressure"),
			"pressure",
			nil,
			nil,
		),
	}

}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up
	ch <- e.version
	ch <- e.temperature
	ch <- e.humidity
	ch <- e.pressure
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	up := float64(1)

	resp, err := http.Get(e.nodeEndpoint)
	if err != nil {
		// notifies that the exporter is down
		ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, 0)
		level.Error(e.logger).Log("msg", "failed to connect to endpoint", "err", err)
		return
	}

	if err := e.parseStats(ch, resp); err != nil {
		level.Error(e.logger).Log("msg", "failed to parse response", "err", err)
		up = 0
	}

	// notifies that the exporter is up
	ch <- prometheus.MustNewConstMetric(e.up, prometheus.GaugeValue, up)
}

func (e *Exporter) parseStats(ch chan<- prometheus.Metric, resp *http.Response) error {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	var metric map[string]interface{}
	if err := json.Unmarshal(body, &metric); err != nil {
		return errors.Wrapf(err, "failed to unmarshal response body")
	}

	// pushes metrics data respectively
	ch <- prometheus.MustNewConstMetric(e.temperature, prometheus.GaugeValue, metric["c"].(float64))
	ch <- prometheus.MustNewConstMetric(e.humidity, prometheus.GaugeValue, metric["h"].(float64))
	ch <- prometheus.MustNewConstMetric(e.pressure, prometheus.GaugeValue, metric["p"].(float64))

	return nil
}
