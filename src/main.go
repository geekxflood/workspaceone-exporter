package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: Add some logs output
// TODO: Set some flags so to be able to select if we want values per tags or not as it can consume tons of API calls
// TODO: Introduce a throttling mechanism to avoid overloading the WS1 API
// TODO: Create a subprocess for quering getting the devices inventory
// TODO: Timeout the API call and produce a metric of this

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
	// fmt.Printf("Type of ws1Intervalraw: %T\n", ws1Intervalraw)
	// fmt.Printf("Value of ws1Intervalraw: %q\n", ws1Intervalraw)

	if ws1Interval, err := strconv.Atoi(ws1Intervalraw); err != nil {
		log.Panic("Error converting WS1_INTERVAL to int")
	} else {
		// fmt.Printf("Type of ws1Interval: %T\n", ws1Interval)
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
			if lastSeen, err := time.Parse("2006-01-02T15:04:05", device.LastSeen); err != nil {
				log.Panic("Error converting lastSeen to time")
			} else {
				// If the lastSeen value - current time is greater than the interval
				// then the device is offline
				if time.Since(lastSeen) > ws1IntervalDuration {
					offline++
				} else {
					online++
				}
			}
		}

		// Set the value of the metric deviceOffline
		deviceOffline.Set(float64(offline))

		// Set the value of the metric deviceOnline
		deviceOnline.Set(float64(online))
	}

	ws1TagParsingRaw := os.Getenv("TAG_PARSING")
	// fmt.Printf("Type of ws1TagParsingRaw: %T\n", ws1TagParsingRaw)
	// fmt.Printf("Value of ws1TagParsingRaw: %q\n", ws1TagParsingRaw)

	if ws1TagParsing, err := strconv.ParseBool(ws1TagParsingRaw); err != nil {
		log.Panic("Error converting WS1_TAG_PARSING to bool")
	} else {
		if ws1TagParsing {

			// Get the number of tags in the WS1 tenant

			tagList := Ws1TagRetriver()
			// fmt.Println(tagList.Total)

			// Set the value of the metric tagSum
			tagSum.Set(float64(tagList.Total))

			// Test if the env Variable TAG_FILTER
			// is not empty, if not filter the tags.TagName
			// that match the TAG_FILTER value
			tagFilter := os.Getenv("TAG_FILTER")
			re := regexp.MustCompile(tagFilter)

			// Now we need to count the number of device per Tag
			// count the number of device per tag
			// For each tag, count the number of device
			// that have this tag

			// Make query to get the list of device per tag
			// But only for the tag that match the TAG_FILTER value

			for _, tag := range tagList.Tags {
				// tag must match the TAG_FILTER value
				if re.MatchString(tag.TagName) {
					tagDeviceList := Ws1TagDeviceRetriver(tag.ID.Value)

					// For each tag found,
					// Set a gauge metric with the tag name and the number of device
					tagDeviceSum.WithLabelValues(tag.TagName).Set(float64(len(tagDeviceList.Device)))

					// For each tag found,
					// Set a gauge metric with the tag name and the number of devices offlline

					for _, tagDevice := range tagDeviceList.Device {

						// Find the device in the deviceList with the same ID
						// and get the lastSeen value
						for _, d := range deviceList.Devices {
							if d.ID.Value == tagDevice.DeviceID {
								// convert the lastSeen value in time
								if lastSeen, err := time.Parse("2006-01-02T15:04:05", d.LastSeen); err != nil {
									panic("Error converting lastSeen to time")
								} else {
									// If the lastSeen value - current time is greater than the interval
									// then the device is offline
									// Get the interval from the environement variable convert the string in int
									ws1Intervalraw := os.Getenv("WS1_INTERVAL")
									// fmt.Printf("Type of ws1Intervalraw: %T\n", ws1Intervalraw)
									// fmt.Printf("Value of ws1Intervalraw: %q\n", ws1Intervalraw)
									if ws1Interval, err := strconv.Atoi(ws1Intervalraw); err != nil {
										panic("Error converting WS1_INTERVAL to int")
									} else {
										if time.Since(lastSeen) > time.Duration(ws1Interval)*time.Minute {
											tagDeviceOffline.WithLabelValues(tag.TagName).Inc()
										}
									}

								}
							}
						}

					}

				}
			}
		}
	}

	// Set the http Handler for prometheus
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9740", nil)
}
