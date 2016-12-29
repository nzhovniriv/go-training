package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	flag.Parse()
	url := flag.Arg(0)
	if len(url) == 0 {
		fmt.Println("Please pass a URL as an argument via the command line")
	} else {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", content)
	}
}
