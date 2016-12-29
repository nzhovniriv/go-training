package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		fmt.Println("Please pass any argument via the command line")
	} else {
		for _, arg := range files {
			file, err := os.Open(arg)
			if err != nil {
				log.Fatal(err)
			}
			countLines(file, arg, counts)
			file.Close()
		}
		for fileName, number := range counts {
			fmt.Printf("%s\t%d\n", fileName, number)
		}
	}
}

// Calculate count lines.
func countLines(file *os.File, arg string, counts map[string]int) {
	input := bufio.NewScanner(file)
	counts[arg] = 0 //if file is empty
	for input.Scan() {
		counts[arg]++
	}
}
