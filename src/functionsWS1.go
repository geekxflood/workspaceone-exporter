package main

// Contains all functions specific to communicate with WorkspaceOne UEM API

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

// function Ws1DeviceRetriver that takes no input, it will return a DevicesResponseObject
// Of all devices in the WS1 tenant
func Ws1DeviceRetriver() DevicesResponseObject {
	// fmt.Println("Retrieving Devices")
	url := os.Getenv("WS1_URL") + "/API/mdm/devices/search?lgid=" + os.Getenv("WS1_LGID")
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
		log.Println("Error: ", err)
	}

	if err != nil {
		log.Println(err)
	}

	// Create the response object
	var responseObject DevicesResponseObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		log.Panic(err)
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
			url := os.Getenv("WS1_URL") + "/API/mdm/devices/search?lgid=" + os.Getenv("WS1_LGID") + "&page=" + strconv.Itoa(i)
			resBody, resStatus, err = ApiCaller(url, method, nil, header)
			//fmt.Println(string(resBody))
			//fmt.Println(resStatus)

			if resStatus != 200 {
				log.Println("Error: ", err)
			}

			if err != nil {
				log.Println(err)
			}

			// Create the response object
			var responseObject2 DevicesResponseObject
			// Unmarshal the response body into the responseObject
			err = json.Unmarshal(resBody, &responseObject2)
			if err != nil {
				log.Panic(err)
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
	url := os.Getenv("WS1_URL") + "/API/mdm/tags/search?organizationgroupid=" + os.Getenv("WS1_LGID")
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
		log.Println("Error: ", err)
	}

	if err != nil {
		log.Panic(err)
	}

	// Create the response object
	var responseObject TagsResponseObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		log.Panic(err)
	}

	return responseObject
}

func Ws1TagDeviceRetriver(tagId int) TagDeviceListObject {
	// fmt.Println("Retrieving Devices")

	id := strconv.Itoa(tagId)

	url := os.Getenv("WS1_URL") + "/API/mdm/tags/" + id + "/devices"
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
		log.Println("Error: ", err)
	}

	if err != nil {
		log.Panic(err)
	}

	// Create the response object
	var responseObject TagDeviceListObject
	// Unmarshal the response body into the responseObject
	err = json.Unmarshal(resBody, &responseObject)
	if err != nil {
		log.Panic(err)
	}

	return responseObject
}
