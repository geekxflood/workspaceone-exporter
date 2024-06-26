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

	ws1TagParsingRaw = os.Getenv("TAG_PARSING")

	tagFilter = os.Getenv("TAG_FILTER")
}

func main() {
	log.Println("Starting WS1 exporter on port 9740")

	// create a new ServeMux
	mux := http.NewServeMux()

	if insecure {
		log.Println("Client insecure:", insecure)
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
									if ws1Interval, err := strconv.Atoi(ws1Intervalraw); err != nil {
										log.Println("Error converting WS1_INTERVAL to int")
									} else {
										// The device is online as it's lastSeen value is less or equal than the interval
										if time.Since(lastSeen) <= time.Duration(ws1Interval)*time.Minute {
											tagDeviceOnline.WithLabelValues(tag.TagName, d.Model).Inc()
											// The device is offline as it's lastSeen value is greater than the interval
										} else if time.Since(lastSeen) > time.Duration(ws1Interval)*time.Minute && time.Since(lastSeen) <= time.Duration(30*24*60)*time.Minute {
											tagDeviceOffline.WithLabelValues(tag.TagName, d.Model).Inc()

											// Test if the time.Since(lastSeen) is greater than 1 month and if so increment the tagDeviceOffline1M metric
										} else if time.Since(lastSeen) > time.Duration(30*24*60)*time.Minute {
											tagDeviceOffline1M.WithLabelValues(tag.TagName, d.Model).Inc()
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
