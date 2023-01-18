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

// TODO: Set some flags so to be able to select if we want values per tags or not as it can consume tons of API calls
// TODO: Introduce a throttling mechanism to avoid overloading the WS1 API
// TODO: Create a subprocess for quering getting the devices inventory
// TODO: Timeout the API call and produce a metric of this

func init() {

	// Test if env variable INSECURE exists
	// If it does, set the insecureSkipVerify to true
	// If it does not, set the insecureSkipVerify to false
	insecure = false
	if os.Getenv("INSECURE") != "" {
		if insecureVal, err := strconv.ParseBool(os.Getenv("INSECURE")); err == nil {
			if insecureVal {
				insecure = true
			}
		}
	}

	// Set the log output to stdout
	log.SetOutput(os.Stdout)
	// Set the log prefix
	log.SetPrefix("ws1_exporter: ")
	// Set the log flag
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Get the interval from the environement variable convert the string in int
	ws1Intervalraw = os.Getenv("WS1_INTERVAL")
	// fmt.Printf("Type of ws1Intervalraw: %T\n", ws1Intervalraw)
	// fmt.Printf("Value of ws1Intervalraw: %q\n", ws1Intervalraw)

	ws1TagParsingRaw = os.Getenv("TAG_PARSING")
	//log.Printf("Type of ws1TagParsingRaw: %T\n", ws1TagParsingRaw)
	//log.Printf("Value of ws1TagParsingRaw: %q\n", ws1TagParsingRaw)$

	tagFilter = os.Getenv("TAG_FILTER")
}

func main() {
	log.Println("Starting WS1 exporter on port 9740")

	// create a new ServeMux
	mux := http.NewServeMux()

	if insecure {
		log.Println("Client insecure: ", insecure)
		SetInsecureSSL()
	}

	deviceList := GetVolumeDevices()

	// List all the diffrent Platform in the inventory
	GetDevicePlatforms(deviceList)

	if ws1Interval, err := strconv.Atoi(ws1Intervalraw); err != nil {
		log.Panic("Error converting WS1_INTERVAL to int")
	} else {

		GetVolumeStatusDevice(deviceList, ws1Interval)
	}

	if ws1TagParsing, err := strconv.ParseBool(ws1TagParsingRaw); err != nil {
		log.Panic("Error converting WS1_TAG_PARSING to bool")
	} else {
		if ws1TagParsing {

			// Get the number of tags in the WS1 tenant

			tagList := Ws1TagRetriver()
			// fmt.Println(tagList.Total)

			GetTagSum(tagList)

			// Test if the env Variable TAG_FILTER
			// is not empty, if not filter the tags.TagName
			// that match the TAG_FILTER value

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
					// Do not set if the value is 0
					if len(tagDeviceList.Device) > 0 {
						tagDeviceSum.WithLabelValues(tag.TagName).Set(float64(len(tagDeviceList.Device)))

					}

					// For each tag found,
					// Set a gauge metric with the tag name and the number of devices offlline

					for _, tagDevice := range tagDeviceList.Device {

						// Find the device in the deviceList with the same ID
						// and get the lastSeen value
						for _, d := range deviceList.Devices {
							if d.ID.Value == tagDevice.DeviceID {
								// convert the lastSeen value in time
								if lastSeen, err := time.Parse("2006-01-02T15:04:05", d.LastSeen); err != nil {
									log.Println("Error converting lastSeen to time")
								} else {
									// If the lastSeen value - current time is greater than the interval
									// then the device is offline
									// Get the interval from the environement variable convert the string in int
									ws1Intervalraw := os.Getenv("WS1_INTERVAL")
									// fmt.Printf("Type of ws1Intervalraw: %T\n", ws1Intervalraw)
									// fmt.Printf("Value of ws1Intervalraw: %q\n", ws1Intervalraw)
									if ws1Interval, err := strconv.Atoi(ws1Intervalraw); err != nil {
										log.Println("Error converting WS1_INTERVAL to int")
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

	// handle GET requests to /metrics
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	// handle all other requests
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	// start the HTTP server
	http.ListenAndServe(":9740", mux)
}
