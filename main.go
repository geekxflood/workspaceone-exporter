package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// TODO: Lookup for moving the types on another file (used to workd I don't know why)
// TODO: Create generic library to host the ApiCaller function, I'll need it for other projects
// TODO: need to lookup what the ... does, from `Code`on discord "It's a destructuring operator, it's like the spread operator in JS"
// TODO: READ THIS -> It's a variadic function https://gobyexample.com/variadic-functions

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
	url := os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("LGID")
	method := "GET"

	header := map[string]string{
		"accept":         "application/json",
		"aw-tenant-code": os.Getenv("WS1_TENANT_KEY"),
		"Authorization":  os.Getenv("WS1_AUTH_KEY"),
		"Content-Type":   "application/json",
	}

	resBody, resStatus, err := ApiCaller(url, method, nil, header)
	//fmt.Println(string(resBody))
	//fmt.Println(resStatus)

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
			url := os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("LGID") + "&page=" + strconv.Itoa(i)
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

func main() {
	deviceList := Ws1DeviceRetriver()
	// Print all the deviceFriendlyName
	for _, device := range deviceList.Devices {
		fmt.Println(device.DeviceFriendlyName)
	}
	// Print the size of the device list
	fmt.Println(len(deviceList.Devices))
}
