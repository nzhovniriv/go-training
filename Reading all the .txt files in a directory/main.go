package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	counts := make(map[string]int)
	flag.Parse()
	dirname := flag.Arg(0)
	if len(dirname) == 0 {
		fmt.Println("Please pass a directory name as an argument via the command line")
	} else {
		dir, err := os.Open(dirname)
		if err != nil {
			log.Fatal(err)
		}
		d, err := dir.Readdir(-1)
		if err != nil {
			log.Fatal(err)
		}
		dir.Close()
		var dirNames []string
		for _, f := range d {
			mode := f.Mode()
			if mode.IsRegular() {
				if match, _ := regexp.MatchString(".txt$", f.Name()); match {
					dirNames = append(dirNames, f.Name())
				}
			}
		}
		readFiles(dirNames, counts, dirname)
		printResult(counts)
	}
}

// Read files in the directory
func readFiles(dirNames []string, counts map[string]int, dirname string) {
	for _, arg := range dirNames {
		file, err := os.Open(dirname + string(filepath.Separator) + arg)
		if err != nil {
			log.Fatal(err)
		}
		countLines(file, arg, counts)
		file.Close()
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

//Print result
func printResult(counts map[string]int) {
	for fileName, number := range counts {
		fmt.Printf("%s\t%d\n", fileName, number)
	}
}
