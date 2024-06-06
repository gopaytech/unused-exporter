package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gopaytech/unused-exporter/pkg/collector"
	"github.com/gopaytech/unused-exporter/pkg/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	settings := settings.NewSettings()
	collectorHandler, err := collector.NewCollector(settings)
	if err != nil {
		log.Fatal(err)
	}

	prometheus.MustRegister(collectorHandler)
	http.Handle("/metrics", promhttp.Handler())

	err = http.ListenAndServe(":"+strconv.Itoa(settings.Port), nil)
	log.Fatal(err)
}
