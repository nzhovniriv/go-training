package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	io.WriteString(w, "My IP address is: "+sendRequest(client))
}

func sendRequest(client *http.Client) string {
	result := make(map[string]string)
	var link url.URL
	link.Scheme = "http"
	link.Host = "httpbin.org"
	link.Path = "ip"
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["origin"]
}
