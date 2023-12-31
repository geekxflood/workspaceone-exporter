// internal/api/api.go

package api

import (
	"io"
	"net/http"

	"github.com/geekxflood/workspaceone-exporter/internal/models"
)

// Example API call function. Extend this with specific API calls as needed.
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

func FetchDevices(client *http.Client, lgid int) (DevicesResponseObject, error) {
	// Implement the logic to fetch devices from the WorkspaceOne UEM API
	// Parse the response and return DevicesResponseObject
}

func FetchTags(client *http.Client, lgid int, tagFilter string) (TagsResponseObject, error) {
	// Implement the logic to fetch tags from the WorkspaceOne UEM API
	// Filter tags if necessary
	// Parse the response and return TagsResponseObject
}
