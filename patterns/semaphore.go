package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	inputJson := []byte(`[
		"http://www.ideaura.com/photos/nature/thumbnails/thumb_ideaura_nature_050.jpg",
		"http://www.soundstacks.co.uk/wp-content/uploads/2016/08/nature-category-250x250.jpeg",
		"http://www.junkinside.com/wp-content/uploads/2009/12/nature-wallpaper-for-mobile-phone-240x320-100x150.jpg",
		"http://thedeershire.com/wp-content/plugins/widgetkit/cache/nature-10-20d75b2ed2.jpg",
		"http://ippcdn1.ippawards.com/wp-content/uploads/2014/10/AndrewVanderWall02-nature-250x250.jpg"
	]`)
	flag.Parse()
	limitation := flag.Arg(0)
	if len(limitation) == 0 {
		fmt.Println("Please pass a limitation value as a flag parameter.")
	} else {
		urls := []string{}
		err := json.Unmarshal(inputJson, &urls)
		if err != nil {
			log.Fatal(err)
		}
		capacity, _ := strconv.Atoi(limitation)
		ch := make(chan struct{}, capacity)
		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				downloadFile(url)
				ch <- struct{}{}
			}(url)
		}
		for range urls {
			<-ch
		}
		wg.Wait()
	}
}

//download file
func downloadFile(url string) {
	fileName := strings.Split(url, "/")
	file, err := os.Create(fileName[len(fileName)-1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}
