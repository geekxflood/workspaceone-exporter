// internal/api/api.go

package api

import (
	"io"
	"net/http"
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
