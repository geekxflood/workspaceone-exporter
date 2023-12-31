package main

// Libraries of all functions that will do any changes on the prometheus metrics

import (
	"log"
	"time"
)

func GetVolumeDevices() (deviceList DevicesResponseObject) {
	// Call the function to get the devices Inventory
	deviceList = Ws1DeviceRetriver()
	deviceNumber.Set(float64(deviceList.Total))
	return deviceList
}

func GetDevicePlatforms(deviceList DevicesResponseObject) {
	// List all the diffrent Platform in the inventory
	platforms := make(map[string]int)
	for _, device := range deviceList.Devices {
		platforms[device.Platform]++
	}
	// For each diffrent platform found, set the value in the metric
	for platform, count := range platforms {
		devicePlatform.WithLabelValues(platform).Set(float64(count))
	}
}

func GetTagSum(tagList TagsResponseObject) {
	// Set the value of the metric tagSum
	tagSum.Set(float64(tagList.Total))
}

func GetVolumeStatusDevice(deviceList DevicesResponseObject, ws1Interval int) {
	// Find the number of device offline
	// The interval is by en evironement variable in minutes
	// This interval is configure in WS1 admin panel and define the time
	// between each check from Ws1 to the device
	// Then we need to find the number of device with a lastSeen
	// older than the interval

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
