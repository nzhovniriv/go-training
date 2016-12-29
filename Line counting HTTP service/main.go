package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Book struct {
	Title string `json:"title"`
	Lines int    `json:"lines,string"`
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var count int
	urlPart := strings.Split(r.URL.Path, "/")
	bookName := urlPart[len(urlPart)-1]
	if len(urlPart) <= 2 || len(bookName) == 0 {
		fmt.Fprint(w, "Please pass all necessary information into URL.")
	} else if len(urlPart) == 3 && len(bookName) != 0 {
		file, err := os.Open(urlPart[len(urlPart)-2] + string(filepath.Separator) + bookName)
		if err != nil {
			log.Fatal(err)
		}
		input := bufio.NewScanner(file)
		for input.Scan() {
			count++
		}
		file.Close()
		book := Book{bookName, count}
		data, err := json.Marshal(book)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	} else {
		fmt.Fprint(w, "You have passed more information into URL.")
	}
}
