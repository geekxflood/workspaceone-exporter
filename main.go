package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func main() {

	recordMetrics()

	// Call the function to get the devices Inventory
	deviceList := Ws1DeviceRetriver()
	deviceNumber.Set(float64(deviceList.Total))

	// Listt all the diffrent Platform in the inventory
	platforms := make(map[string]int)
	for _, device := range deviceList.Devices {
		platforms[device.Platform]++
	}
	// For each diffrent platform found, set the value in the metric
	for platform, count := range platforms {
		devicePlatform.WithLabelValues(platform).Set(float64(count))
	}

	// Find the number of device offline
	// The interval is by en evironement variable in minutes
	// This interval is configure in WS1 admin panel and define the time
	// between each check from Ws1 to the device
	// Then we need to find the number of device with a lastSeen
	// older than the interval

	// Get the interval from the environement variable convert the string in int
	ws1Intervalraw := os.Getenv("WS1_INTERVAL")
	fmt.Printf("Type of ws1Intervalraw: %T\n", ws1Intervalraw)
	fmt.Printf("Value of ws1Intervalraw: %q\n", ws1Intervalraw)
	ws1Interval, err := strconv.Atoi(ws1Intervalraw)
	if err != nil {
		panic("Error converting WS1_INTERVAL to int")
	}
	fmt.Printf("Type of ws1Interval: %T\n", ws1Interval)

	// Convert ws1Interval value in minutes time.duration
	ws1IntervalDuration := time.Duration(ws1Interval) * time.Minute
	// fmt.Printf("Type of ws1IntervalDuration: %T\n", ws1IntervalDuration)
	// Print the value of ws1IntervalDuration
	// fmt.Printf("Value of ws1IntervalDuration: %q\n", ws1IntervalDuration)

	// For each device, evaluate the lastSeen value
	// If the lastSeen value - current time is greater than the interval
	// then the device is offline
	offline := 0
	online := 0
	for _, device := range deviceList.Devices {
		// convert the lastSeen value in time
		lastSeen, err := time.Parse("2006-01-02T15:04:05", device.LastSeen)
		if err != nil {
			panic(err)
		}
		// If the lastSeen value - current time is greater than the interval
		// then the device is offline
		if time.Since(lastSeen) > ws1IntervalDuration {
			offline++
		} else {
			online++
		}
	}

	// Set the value of the metric deviceOffline
	deviceOffline.Set(float64(offline))

	// Set the value of the metric deviceOnline
	deviceOnline.Set(float64(online))

	// Get the number of tags in the WS1 tenant

	tagList := Ws1TagRetriver()
	_ = tagList
	// fmt.Println(tagList.Total)

	// Set the http Handler for prometheus
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9740", nil)
}
