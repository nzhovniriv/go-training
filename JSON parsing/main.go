package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := "http://httpbin.org/ip"
	response, ok := http.Get(url)
	if ok != nil {
		log.Fatal(ok)
	}
	result := make(map[string]string)
	err := json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	fmt.Println("My IP address is: " + result["origin"])
}
