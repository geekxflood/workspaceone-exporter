package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// TODO: Create a function that will count the number of devices per TAG
// TODO: Create a function that will count the number of devices offline
// TODO: offline status can be define as a device that has it's LastSeen value older than refresh server time
// TODO: Find a way to query the refresh server time so to not hard code it
// TODO: otherwise we will specify it using a env variable that will represent a value in minutes

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

// function Ws1TagRetriver will return a TagsResponseObject
// list of all tags in the WS1 tenant
func Ws1TagRetriver() TagsResponseObject {
	url := os.Getenv("WS1_URL") + "/mdm/tags/search?organizationgroupid=" + os.Getenv("LGID")
	method := "GET"

	header := map[string]string{
		"accept":         "application/json",
		"aw-tenant-code": os.Getenv("WS1_TENANT_KEY"),
		"Authorization":  os.Getenv("WS1_AUTH_KEY"),
		"Content-Type":   "application/json",
	}

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
