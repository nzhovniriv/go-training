package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	inputJson := []byte(`[
		"http://www.ideaura.com/photos/nature/thumbnails/thumb_ideaura_nature_050.jpg",
		"http://www.soundstacks.co.uk/wp-content/uploads/2016/08/nature-category-250x250.jpeg",
		"http://www.junkinside.com/wp-content/uploads/2009/12/nature-wallpaper-for-mobile-phone-240x320-100x150.jpg",
		"http://thedeershire.com/wp-content/plugins/widgetkit/cache/nature-10-20d75b2ed2.jpg",
		"http://ippcdn1.ippawards.com/wp-content/uploads/2014/10/AndrewVanderWall02-nature-250x250.jpg"
	]`)
	urls := []string{}
	err := json.Unmarshal(inputJson, &urls)
	if err != nil {
		log.Fatal(err)
	}
	result := findFirstResult(urls)
	fmt.Println(result)
}

func findFirstResult(urls []string) string {
	ch := make(chan string, 1)
	for _, url := range urls {
		go func(url string) {
			select {
			case ch <- request(url):
			default:
			}
		}(url)
	}
	return <-ch
}

func request(url string) string {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return url
}
