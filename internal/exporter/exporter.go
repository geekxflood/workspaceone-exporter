package exporter

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/geekxflood/workspaceone-exporter/internal/api"
	"github.com/geekxflood/workspaceone-exporter/internal/metrics"
	"github.com/geekxflood/workspaceone-exporter/internal/models"
)

func FetchAndUpdateMetrics(client *http.Client, lgid int, ws1Interval string, tagParsing bool, tagFilter string) {
	// Fetch device list
	deviceList, err := api.FetchDevices(client, lgid)
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	metrics.DeviceNumber.Set(float64(len(deviceList.Devices)))

	// Process and update metrics for devices
	processAndUpdateDeviceMetrics(deviceList, ws1Interval)

	if tagParsing {
		// Fetch and process tags if tag parsing is enabled
		tagList, err := api.FetchTags(client, lgid, tagFilter)
		if err != nil {
			log.Printf("Error fetching tags: %v", err)
			return
		}

		metrics.TagSum.Set(float64(len(tagList.Tags)))

		// Process and update metrics for tags
		processAndUpdateTagMetrics(client, tagList, deviceList, ws1Interval)
	}
}

func processAndUpdateDeviceMetrics(deviceList api.DevicesResponseObject, ws1Interval string) {
	interval, err := strconv.Atoi(ws1Interval)
	if err != nil {
		log.Printf("Error converting WS1_INTERVAL to int: %v", err)
		return
	}
	intervalDuration := time.Duration(interval) * time.Minute

	var offlineCount, onlineCount int
	platforms := make(map[string]int)
	for _, device := range deviceList.Devices {
		platforms[device.Platform]++
		lastSeen, err := time.Parse("2006-01-02T15:04:05", device.LastSeen)
		if err != nil {
			log.Printf("Error parsing lastSeen: %v", err)
			continue
		}
		if time.Since(lastSeen) > intervalDuration {
			offlineCount++
		} else {
			onlineCount++
		}
	}

	metrics.DeviceOffline.Set(float64(offlineCount))
	metrics.DeviceOnline.Set(float64(onlineCount))

	for platform, count := range platforms {
		metrics.DevicePlatform.WithLabelValues(platform).Set(float64(count))
	}
}

func processAndUpdateTagMetrics(client *http.Client, tagList api.TagsResponseObject, deviceList api.DevicesResponseObject, ws1Interval string) {
	// Implement logic to process and update metrics for tags
	// Similar to device metrics, but with additional logic for handling tags
}
