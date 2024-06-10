package model

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IPAddress struct {
	Cloud    string
	Region   string
	Identity string
	Value    string
	Type     string
	Used     bool
}

type LoadBalancer struct {
	Cloud  string
	Region string
	IP     string
	Type   string
	Used   bool
}

var (
	IPAddressUnusedGauge = prometheus.NewDesc(
		prometheus.BuildFQName("ip_address", "", "unused"),
		"Indicates unused IP Address",
		[]string{"cloud", "region", "value", "type", "identity"}, nil,
	)

	IPAddressUsedGauge = prometheus.NewDesc(
		prometheus.BuildFQName("ip_address", "", "used"),
		"Indicates used IP Address",
		[]string{"cloud", "region", "value", "type", "identity"}, nil,
	)

	LoadBalancerUnusedGauge = prometheus.NewDesc(
		prometheus.BuildFQName("load_balancer", "", "unused"),
		"Indicates unused Load Balancer",
		[]string{"cloud", "region", "ip", "type"}, nil,
	)

	LoadBalancerUsedGauge = prometheus.NewDesc(
		prometheus.BuildFQName("load_balancer", "", "used"),
		"Indicates used Load Balancer",
		[]string{"cloud", "region", "ip", "type"}, nil,
	)
)
