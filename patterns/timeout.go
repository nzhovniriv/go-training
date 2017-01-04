package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()
	content, err := request(ts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(content)
}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Fprintln(w, "Hello, world!")
}

func request(ts *httptest.Server) (string, error) {
	response, err := http.Get(ts.URL)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil
}
