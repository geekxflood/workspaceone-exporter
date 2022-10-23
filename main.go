package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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
	fmt.Printf("%T\n", body)

	//fmt.Println(string(body))
}
