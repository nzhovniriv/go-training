package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Book struct {
	Title string `json:"title"`
	Lines int32  `json:"lines,string"`
}

func main() {
	url := "http://127.0.0.1:3999/books/"
	var count int32
	flag.Parse()
	title := flag.Arg(0)
	if len(title) == 0 {
		fmt.Println("Please pass a book name as an argument via the command line")
	} else {
		response, err := http.Get(url + title)
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(response.Body)
		for {
			_, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break
			} else {
				count++
			}
		}
		response.Body.Close()
		book := Book{title, count}
		data, err := json.Marshal(book)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", data)
	}
}
