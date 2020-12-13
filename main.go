package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	up = prometheus.NewDesc(
		"envii_up",
		"was talking to envii successful.",
		nil, nil,
	)
)

type EnviiCollector struct {
}

type EnviiMetric struct {
	C float64 `json:"c"`
	H float64 `json:"h"`
	P float64 `json:"p"`
}

func (e EnviiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (e EnviiCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get("http://192.168.11.12:9000")
	if err != nil {
		log.Println("failed to get response")
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("failed to read response body")
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}

	var metric map[string]interface{}
	if err := json.Unmarshal(body, &metric); err != nil {
		log.Println("failed to unmarshal response body")
		ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc("envii_temperature", "temperature", nil, nil),
		prometheus.GaugeValue, metric["c"].(float64))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc("envii_humidity", "humidity", nil, nil),
		prometheus.GaugeValue, metric["h"].(float64))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc("envii_pressure", "air_pressure", nil, nil),
		prometheus.GaugeValue, metric["p"].(float64))
}

func main() {
	e := EnviiCollector{}
	prometheus.MustRegister(e)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:8001", nil))
}
