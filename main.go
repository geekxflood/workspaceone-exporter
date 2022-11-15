package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TODO: Add some logs output
// TODO: Set some flags so to be able to select if we want values per tags or not as it can consume tons of API calls
// TODO: Introduce a throttling mechanism to avoid overloading the WS1 API
// TODO: Create a subprocess for quering getting the devices inventory
// TODO: Timeout the API call and produce a metric of this

type TagsResponseObject struct {
	Tags []struct {
		TagName    string `json:"TagName"`
		DateTagged string `json:"DateTagged"`
		TagAvatar  string `json:"TagAvatar"`
		ID         struct {
			Value int `json:"Value"`
		} `json:"Id"`
		UUID string `json:"Uuid"`
	} `json:"Tags"`
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
	Total    int `json:"Total"`
}

type TagDeviceListObject struct {
	Device []struct {
		DeviceID     int    `json:"DeviceId"`
		FriendlyName string `json:"FriendlyName"`
		DateTagged   string `json:"DateTagged"`
		DeviceUUID   string `json:"DeviceUuid"`
	} `json:"Device"`
}

type TagInvDeviceObject struct {
	TagName string
	TagID   int
	Device  []struct {
		DeviceID     int
		FriendlyName string
	}
}

type DevicesResponseObject struct {
	Devices []struct {
		EasIds struct {
		} `json:"EasIds"`
		TimeZone           string `json:"TimeZone"`
		Udid               string `json:"Udid"`
		SerialNumber       string `json:"SerialNumber"`
		MacAddress         string `json:"MacAddress"`
		Imei               string `json:"Imei"`
		EasID              string `json:"EasId"`
		AssetNumber        string `json:"AssetNumber"`
		DeviceFriendlyName string `json:"DeviceFriendlyName"`
		DeviceReportedName string `json:"DeviceReportedName"`
		LocationGroupID    struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
			UUID string `json:"Uuid"`
		} `json:"LocationGroupId"`
		LocationGroupName string `json:"LocationGroupName"`
		UserID            struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
			UUID string `json:"Uuid"`
		} `json:"UserId"`
		UserName             string `json:"UserName"`
		DataProtectionStatus int    `json:"DataProtectionStatus"`
		UserEmailAddress     string `json:"UserEmailAddress"`
		Ownership            string `json:"Ownership"`
		PlatformID           struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
		} `json:"PlatformId"`
		Platform string `json:"Platform"`
		ModelID  struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
		} `json:"ModelId"`
		Model                            string `json:"Model"`
		OperatingSystem                  string `json:"OperatingSystem"`
		PhoneNumber                      string `json:"PhoneNumber"`
		LastSeen                         string `json:"LastSeen"`
		EnrollmentStatus                 string `json:"EnrollmentStatus"`
		ComplianceStatus                 string `json:"ComplianceStatus"`
		CompromisedStatus                bool   `json:"CompromisedStatus"`
		LastEnrolledOn                   string `json:"LastEnrolledOn"`
		LastComplianceCheckOn            string `json:"LastComplianceCheckOn"`
		LastCompromisedCheckOn           string `json:"LastCompromisedCheckOn"`
		IsSupervised                     bool   `json:"IsSupervised"`
		VirtualMemory                    int    `json:"VirtualMemory"`
		OEMInfo                          string `json:"OEMInfo"`
		IsDeviceDNDEnabled               bool   `json:"IsDeviceDNDEnabled"`
		IsDeviceLocatorEnabled           bool   `json:"IsDeviceLocatorEnabled"`
		IsCloudBackupEnabled             bool   `json:"IsCloudBackupEnabled"`
		IsActivationLockEnabled          bool   `json:"IsActivationLockEnabled"`
		IsNetworkTethered                bool   `json:"IsNetworkTethered"`
		BatteryLevel                     string `json:"BatteryLevel"`
		IsRoaming                        bool   `json:"IsRoaming"`
		SystemIntegrityProtectionEnabled bool   `json:"SystemIntegrityProtectionEnabled"`
		ProcessorArchitecture            int    `json:"ProcessorArchitecture"`
		TotalPhysicalMemory              int    `json:"TotalPhysicalMemory"`
		AvailablePhysicalMemory          int    `json:"AvailablePhysicalMemory"`
		OSBuildVersion                   string `json:"OSBuildVersion"`
		DeviceCellularNetworkInfo        []struct {
			CarrierName string `json:"CarrierName"`
			CardID      string `json:"CardId"`
			PhoneNumber string `json:"PhoneNumber"`
			DeviceMCC   struct {
				Simmcc     string `json:"SIMMCC"`
				CurrentMCC string `json:"CurrentMCC"`
			} `json:"DeviceMCC"`
			IsRoaming bool `json:"IsRoaming"`
		} `json:"DeviceCellularNetworkInfo,omitempty"`
		EnrollmentUserUUID string `json:"EnrollmentUserUuid"`
		ManagedBy          int    `json:"ManagedBy"`
		WifiSsid           string `json:"WifiSsid"`
		ID                 struct {
			Value int `json:"Value"`
		} `json:"Id"`
		UUID              string `json:"Uuid"`
		ComplianceSummary struct {
			DeviceCompliance []struct {
				CompliantStatus     bool          `json:"CompliantStatus"`
				PolicyName          string        `json:"PolicyName"`
				PolicyDetail        string        `json:"PolicyDetail"`
				LastComplianceCheck string        `json:"LastComplianceCheck"`
				NextComplianceCheck string        `json:"NextComplianceCheck"`
				ActionTaken         []interface{} `json:"ActionTaken"`
				ID                  struct {
					Value int `json:"Value"`
				} `json:"Id"`
				UUID string `json:"Uuid"`
			} `json:"DeviceCompliance"`
		} `json:"ComplianceSummary,omitempty"`
	} `json:"Devices"`
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
	Total    int `json:"Total"`
}

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

// function apiCaller will do a REST API Call
// The function takes as argument the URL of the API, the method (GET, POST, PUT, DELETE),
// the body of the request (if any) and the headers (if any)
// The function return the response body, the status code and an error
func ApiCaller(url string, method string, body io.Reader, headers map[string]string) ([]byte, int, error) {
	// Create a new request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	// Add the headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Do the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return respBody, resp.StatusCode, nil
}

// function Ws1DeviceRetriver that takes no input, it will return a DevicesResponseObject
// Of all devices in the WS1 tenant
func Ws1DeviceRetriver() DevicesResponseObject {
	// fmt.Println("Retrieving Devices")
	url := os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("WS1_LGID")
	method := "GET"

	header := map[string]string{
		"accept":         "application/json",
		"aw-tenant-code": os.Getenv("WS1_TENANT_KEY"),
		"Authorization":  os.Getenv("WS1_AUTH_KEY"),
		"Content-Type":   "application/json",
	}

	// fmt.Println("Calling API")
	resBody, resStatus, err := ApiCaller(url, method, nil, header)

	if resStatus != 200 {
		fmt.Println("Error: ", err)
	}

	if err != nil {
		fmt.Println(err)
	}

	// Create the response object
	var responseObject DevicesResponseObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		panic(err)
	}

	// Check is the number of device is greater than the page size
	// If it is, we need to iterate on the pages and add the devices to the responseObject.Devices
	if responseObject.Total > responseObject.PageSize {
		// Find the number of pages
		var pages int = responseObject.Total / responseObject.PageSize
		if responseObject.Total%responseObject.PageSize > 0 {
			pages++
		}
		// redo the API call for each page
		// Start at 1 because the first page @ 0 is already in the responseObject
		for i := 1; i < pages; i++ {
			url := os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("WS1_LGID") + "&page=" + strconv.Itoa(i)
			resBody, resStatus, err = ApiCaller(url, method, nil, header)
			//fmt.Println(string(resBody))
			//fmt.Println(resStatus)

			if resStatus != 200 {
				fmt.Println("Error: ", err)
			}

			if err != nil {
				fmt.Println(err)
			}

			// Create the response object
			var responseObject2 DevicesResponseObject
			// Unmarshal the response body into the responseObject
			err = json.Unmarshal(resBody, &responseObject2)
			if err != nil {
				panic(err)
			}
			// Add the devices to the responseObject.Devices
			responseObject.Devices = append(responseObject.Devices, responseObject2.Devices...)
		}
	}

	return responseObject
}

// function Ws1TagRetriver will return a TagsResponseObject
// list of all tags in the WS1 tenant
func Ws1TagRetriver() TagsResponseObject {
	// fmt.Println("Retrieving Tags")
	url := os.Getenv("WS1_URL") + "/mdm/tags/search?organizationgroupid=" + os.Getenv("WS1_LGID")
	method := "GET"

	header := map[string]string{
		"accept":         "application/json",
		"aw-tenant-code": os.Getenv("WS1_TENANT_KEY"),
		"Authorization":  os.Getenv("WS1_AUTH_KEY"),
		"Content-Type":   "application/json",
	}

	// fmt.Println("Calling API")
	resBody, resStatus, err := ApiCaller(url, method, nil, header)

	if resStatus != 200 {
		fmt.Println("Error: ", err)
	}

	if err != nil {
		panic(err)
	}

	// Create the response object
	var responseObject TagsResponseObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		panic(err)
	}

	return responseObject
}

func Ws1TagDeviceRetriver(tagId int) TagDeviceListObject {
	// fmt.Println("Retrieving Devices")

	id := strconv.Itoa(tagId)

	url := os.Getenv("WS1_URL") + "/mdm/tags/" + id + "/devices"
	method := "GET"

	header := map[string]string{
		"accept":         "application/json",
		"aw-tenant-code": os.Getenv("WS1_TENANT_KEY"),
		"Authorization":  os.Getenv("WS1_AUTH_KEY"),
		"Content-Type":   "application/json",
	}

	// fmt.Println("Calling API")
	resBody, resStatus, err := ApiCaller(url, method, nil, header)

	if resStatus != 200 {
		fmt.Println("Error: ", err)
	}

	if err != nil {
		panic(err)
	}

	// Create the response object
	var responseObject TagDeviceListObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		panic(err)
	}

	return responseObject
}

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
		panic("Error converting WS1_INTERVAL to int")
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
				panic("Error converting lastSeen to time")
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

		panic("Error converting WS1_TAG_PARSING to bool")
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
