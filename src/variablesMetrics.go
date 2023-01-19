package main

// Contains all the variables that will be used to create the metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// deviceNumber is a gauge which represents the number of devices in the WS1 tenant
var (
	deviceNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "devices_number",
		Help: "The number of devices in the WS1 tenant",
	})
)

// devicePlatform is a gauge which represents the number of devices per OS in the WS1 tenant
var (
	devicePlatform = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "devices_os",
		Help: "The number of devices per OS in the WS1 tenant",
	}, []string{"platform"})
)

// deviceOffline is a gauge which represents the number of devices in the WS1 tenant that are offline
var (
	deviceOffline = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "devices_offline",
		Help: "The number of devices in the WS1 tenant that are offline",
	})
)

// deviceOnline is a gauge which represents the number of devices in the WS1 tenant that are online
var (
	deviceOnline = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "devices_online",
		Help: "The number of devices in the WS1 tenant that are online",
	})
)

// tagSum is a gauge which represents the number of tags in the WS1 tenant
var (
	tagSum = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tags_sum",
		Help: "The number of tags in the WS1 tenant",
	})
)

// tagDeviceOnline is a gauge wich represent the number of device online per tag in the WS1 tenant
var (
	tagDeviceOnline = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "devices_online_tag",
		Help: "The number of devices online per tag in the WS1 tenant",
	}, []string{"tag", "model"})
)

// tagDeviceOffline is a gauge wich represent the number of device offline per tag in the WS1 tenant
var (
	tagDeviceOffline = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "devices_offline_tag",
		Help: "The number of devices offline per tag in the WS1 tenant",
	}, []string{"tag", "model"})
)

// tagDeviceOffline1M is a gauge wich represent the number of device offline per tag in the WS1 tenant for more than the last month
var (
	tagDeviceOffline1M = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "devices_offline_1m_tag",
		Help: "The number of devices offline per tag in the WS1 tenant for more than the last month",
	}, []string{"tag", "model"})
)

// apiCalls is a counter which represents the number of API calls made to the WS1 API
var (
	apiCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "api_calls",
		Help: "The number of API calls made to the WS1 API",
	})
)
