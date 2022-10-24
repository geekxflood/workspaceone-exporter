package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: Create a subprocess for quering getting the devices inventory
// TODO: Timeout the API call and produce a metric of this
// TODO: need to lookup what the ... does, from `Code`on discord "It's a destructuring operator, it's like the spread operator in JS"
// TODO: READ THIS -> It's a variadic function https://gobyexample.com/variadic-functions

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

// opsProcessed is a counter which represents the total number of processed events.
// This counter is incremented each 2 seconds by the recordMetrics function
var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

// deviceNumber is a gauge which represents the number of devices in the WS1 tenant
// This gauge is incremented each time the Ws1DeviceRetriver function is called
var (
	deviceNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "device_number",
		Help: "The number of devices in the WS1 tenant",
	})
)

func main() {

	recordMetrics()

	// Call the function to get the devices Inventory
	deviceList := Ws1DeviceRetriver()
	deviceNumber.Set(float64(deviceList.Total))

	// Set the http Handler for prometheus
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9740", nil)
}
