package main

import (
	"io"
	"net/http"
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
