package collector

import (
	"github.com/gopaytech/unused-exporter/pkg/data"
	"github.com/gopaytech/unused-exporter/pkg/model"
	"github.com/gopaytech/unused-exporter/pkg/settings"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	settings settings.Settings
	google   data.Data
	aws      data.Data
}

// Describe implements the prometheus.Collector interface.
func (e *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- model.IPAddressUnusedGauge
}

// Collect implements the prometheus.Collector interface.
func (e *Collector) Collect(ch chan<- prometheus.Metric) {
	if e.settings.EnableGCP {
		unusedIPs, err := e.google.GetUnusedIP()
		if err != nil {
			return
		}

		for _, unusedIP := range unusedIPs {
			ch <- prometheus.MustNewConstMetric(model.IPAddressUnusedGauge, prometheus.GaugeValue, 1, unusedIP.Cloud, unusedIP.Region, unusedIP.Value, unusedIP.Type, unusedIP.Identity)
		}

		usedIPs, err := e.google.GetUsedIP()
		if err != nil {
			return
		}

		for _, usedIP := range usedIPs {
			ch <- prometheus.MustNewConstMetric(model.IPAddressUsedGauge, prometheus.GaugeValue, 1, usedIP.Cloud, usedIP.Region, usedIP.Value, usedIP.Type, usedIP.Identity)
		}
	}

	if e.settings.EnableAWS {
		unusedIPs, err := e.aws.GetUnusedIP()
		if err != nil {
			return
		}

		for _, unusedIP := range unusedIPs {
			ch <- prometheus.MustNewConstMetric(model.IPAddressUnusedGauge, prometheus.GaugeValue, 1, unusedIP.Cloud, unusedIP.Region, unusedIP.Value, unusedIP.Type, unusedIP.Identity)
		}

		usedIPs, err := e.aws.GetUsedIP()
		if err != nil {
			return
		}

		for _, usedIP := range usedIPs {
			ch <- prometheus.MustNewConstMetric(model.IPAddressUsedGauge, prometheus.GaugeValue, 1, usedIP.Cloud, usedIP.Region, usedIP.Value, usedIP.Type, usedIP.Identity)
		}
	}
}

func NewCollector(s settings.Settings) (*Collector, error) {
	google, err := data.NewGoogleData(s)
	if err != nil {
		return nil, err
	}

	return &Collector{
		settings: s,
		google:   google,
	}, nil
}
