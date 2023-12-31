// internal/api/api.go

package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/geekxflood/workspaceone-exporter/internal/models"
	"github.com/geekxflood/workspaceone-exporter/internal/util"
)

// API call function
func CallAPI(client *http.Client, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return client.Do(req)
}

// Retrive the devices from the API and return a DevicesResponseObject
func FetchDevices(client *http.Client, lgid string) (models.DevicesResponseObject, error) {
	// Retrieve the necessary environement variables
	envs, err := util.GetEnv([]string{"WS1_AUTH_KEY", "WS1_TENANT_KEY", "WS1_URL"})
	if err != nil {
		return models.DevicesResponseObject{}, err
	}

	// Build the request URL
	url := envs["WS1_URL"] + "/api/mdm/devices/search?pagesize=500&lgid=" + lgid

	// Build the request headers
	headers := map[string]string{
		"aw-tenant-code": envs["WS1_TENANT_KEY"],
		"Authorization":  "Basic " + envs["WS1_AUTH_KEY"],
		"Accept":         "application/json",
	}

	// Make the API call
	resp, err := CallAPI(client, "GET", url, nil, headers)
	if err != nil {
		return models.DevicesResponseObject{}, err
	}
	if resp.StatusCode != 200 {
		return models.DevicesResponseObject{}, err
	}

	// Parse the response body
	var devices models.DevicesResponseObject
	if err := util.ParseBody(resp.Body, &devices); err != nil {
		return models.DevicesResponseObject{}, err
	}

	// Check if there are more devices than the pages size
	// If so, we need to make more API calls to get all the devices
	if devices.Total > devices.Page {
		// Calculate the number of pages to fetch
		pages := devices.Total / devices.Page
		if devices.Total%devices.Page > 0 {
			pages++
		}

		// Make the API calls
		for i := 2; i <= pages; i++ {
			// Build the request URL
			url := envs["WS1_URL"] + "/api/mdm/devices/search?pagesize=500&lgid=" + lgid + "&page=" + fmt.Sprint(i)

			// Make the API call
			resp, err := CallAPI(client, "GET", url, nil, headers)
			if err != nil {
				return models.DevicesResponseObject{}, err
			}
			if resp.StatusCode != 200 {
				return models.DevicesResponseObject{}, err
			}

			// Parse the response body
			var devicesPage models.DevicesResponseObject
			if err := util.ParseBody(resp.Body, &devicesPage); err != nil {
				return models.DevicesResponseObject{}, err
			}

			// Append the devices to the devices slice
			devices.Devices = append(devices.Devices, devicesPage.Devices...)
		}
	}

	return devices, nil

}
