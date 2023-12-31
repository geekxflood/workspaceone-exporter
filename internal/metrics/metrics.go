// internal/metrics/metrics.go

package metrics

import "github.com/prometheus/client_golang/prometheus"

// deviceNumber is a gauge that represents the number of devices in the WS1 tenant.
var DeviceNumber = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "devices_number",
	Help: "The number of devices in the WS1 tenant",
})

// DevicePlatform is a gauge that represents the number of devices per OS in the WS1 tenant.
var DevicePlatform = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "devices_os",
	Help: "The number of devices per OS in the WS1 tenant",
}, []string{"platform"})

// DeviceOffline is a gauge that represents the number of devices in the WS1 tenant that are offline.
var DeviceOffline = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "devices_offline",
	Help: "The number of devices in the WS1 tenant that are offline",
})

// DeviceOnline is a gauge that represents the number of devices in the WS1 tenant that are online.
var DeviceOnline = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "devices_online",
	Help: "The number of devices in the WS1 tenant that are online",
})

// TagSum is a gauge that represents the number of tags in the WS1 tenant.
var TagSum = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "tags_sum",
	Help: "The number of tags in the WS1 tenant",
})

// TagDeviceOnline is a gauge that represents the number of devices online per tag in the WS1 tenant.
var TagDeviceOnline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "devices_online_tag",
	Help: "The number of devices online per tag in the WS1 tenant",
}, []string{"tag", "model"})

// TagDeviceOffline is a gauge that represents the number of devices offline per tag in the WS1 tenant.
var TagDeviceOffline = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "devices_offline_tag",
	Help: "The number of devices offline per tag in the WS1 tenant",
}, []string{"tag", "model"})

// TagDeviceOffline1M is a gauge that represents the number of devices offline per tag in the WS1 tenant for more than the last month.
var TagDeviceOffline1M = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "devices_offline_1m_tag",
	Help: "The number of devices offline per tag in the WS1 tenant for more than the last month",
}, []string{"tag", "model"})

// ApiCalls is a counter that represents the number of API calls made to the WS1 API.
var ApiCalls = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "api_calls",
	Help: "The number of API calls made to the WS1 API",
})

func init() {
	// Register all metrics here.
	prometheus.MustRegister(DeviceNumber, DevicePlatform, DeviceOffline, DeviceOnline, TagSum, TagDeviceOnline, TagDeviceOffline, TagDeviceOffline1M, ApiCalls)
}
