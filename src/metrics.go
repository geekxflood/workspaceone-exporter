package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// opsProcessed is a counter which represents the total number of processed events.
// This counter is incremented each 2 seconds by the recordMetrics function
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

// deviceNumber is a gauge which represents the number of devices in the WS1 tenant
var (
	deviceNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "device_number",
		Help: "The number of devices in the WS1 tenant",
	})
)

// devicePlatform is a gauge which represents the number of devices per OS in the WS1 tenant
var (
	devicePlatform = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "device_os",
		Help: "The number of devices per OS in the WS1 tenant",
	}, []string{"platform"})
)

// deviceOffline is a gauge which represents the number of devices in the WS1 tenant that are offline
var (
	deviceOffline = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "device_offline",
		Help: "The number of devices in the WS1 tenant that are offline",
	})
)

// deviceOnline is a gauge which represents the number of devices in the WS1 tenant that are online
var (
	deviceOnline = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "device_online",
		Help: "The number of devices in the WS1 tenant that are online",
	})
)

// tagSum is a gauge which represents the number of tags in the WS1 tenant
var (
	tagSum = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tag_sum",
		Help: "The number of tags in the WS1 tenant",
	})
)

// tagDeviceSum is a gauge which represents the number of devices per tag in the WS1 tenant
var (
	tagDeviceSum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tag_device_sum",
		Help: "The number of devices per tag in the WS1 tenant",
	}, []string{"tag"})
)

// tagDeviceOffline is a gauge wich represent the number of device offline per tag in the WS1 tenant
var (
	tagDeviceOffline = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "tag_device_offline",
		Help: "The number of device offline per tag in the WS1 tenant",
	}, []string{"tag"})
)
