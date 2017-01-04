package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	arg := flag.Arg(0)
	if len(arg) == 0 {
		fmt.Println("Please pass a number of concurrent servers as a flag parameter.")
	} else {
		count, _ := strconv.Atoi(arg)
		ch := make(chan string)
		for i := 0; i < count; i++ {
			go func() {
				ts := httptest.NewServer(http.HandlerFunc(handler))
				defer ts.Close()
				request(ts, ch)
			}()
		}
		printData(ch)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Fprintln(w, "Hello, world!")
}

func request(ts *httptest.Server, ch chan<- string) {
	for {
		time.Sleep(1 * time.Second)
		response, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		ch <- string(content)
		close(ch)
	}
}

func printData(ch <-chan string) {
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf(data)
	}
}