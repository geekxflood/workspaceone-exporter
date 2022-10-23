package main

// TODO: Do a function for all the API Call it's getting redungdant
// TODO: DO a function that will iterate on the pages and append the devices to the responseObject.Devices
// TODO: need to lookup what the ... does, from `Code`on discord "It's a destructuring operator, it's like the spread operator in JS"
// TODO: READ THIS -> It's a variadic function https://gobyexample.com/variadic-functions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func runHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

}

func main() {
	// Query WS1 service

	url := os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("LGID")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("aw-tenant-code", os.Getenv("WS1_TENANT_KEY"))
	req.Header.Add("Authorization", os.Getenv("WS1_AUTH_KEY"))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Client Error")
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO read Error")
		return
	}
	// fmt.Printf("%T\n", body)

	// convert the body to a DevicesResponseObject
	var responseObject DevicesResponseObject
	json.Unmarshal(body, &responseObject)
	fmt.Println("Number of devices are:", responseObject.Total)
	fmt.Println("Size of page is :", responseObject.PageSize)
	fmt.Println("Number of devices in current page are:", len(responseObject.Devices))

	url = os.Getenv("WS1_URL") + "/mdm/devices/search?lgid=" + os.Getenv("LGID") + "&page=1"
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("aw-tenant-code", os.Getenv("WS1_TENANT_KEY"))
	req.Header.Add("Authorization", os.Getenv("WS1_AUTH_KEY"))

	res, err = client.Do(req)
	if err != nil {
		fmt.Println("Client Error")
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO read Error")
		return
	}

	//fmt.Println(string(body))

	var tmpBody DevicesResponseObject
	json.Unmarshal(body, &tmpBody)

	fmt.Println("Number of devices are:", tmpBody.Total)

	// responseObject.Devices + body.Devices
	responseObject.Devices = append(responseObject.Devices, tmpBody.Devices...)

	fmt.Println("Number of devices after appending are:", len(responseObject.Devices))

	// Loop through the devices and print the device name
	// for _, device := range responseObject.Devices {
	// 	fmt.Println(device.SerialNumber)
	// }

}
