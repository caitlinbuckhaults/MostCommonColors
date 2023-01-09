package main

import (
	"MostCommonColors/imageprocessor"
	"bufio"
	"fmt"
	"os"
	"sync"
)

const (
	maxGoroutines   = 100
	maxHTTPRequests = 1000
)

func main() {
	//set up a buffered channel to pass URLs to the goroutines
	urlChan := make(chan string, maxGoroutines)
	//read input file & parse image URLs
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Problem opening input file: ", err)
	}
	defer file.Close()

	fmt.Println("file successfully opened")
	scanner := bufio.NewScanner(file)

	//foreach url, process image
	var wg sync.WaitGroup
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlChan {
				fmt.Println("Proccessing file at: ", url)
				err := imageprocessor.ProcessImageOptimized(url)
				if err != nil {
					wg.Done()
					return
				}
			}
		}()
	}
	for scanner.Scan() {
		urlChan <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Problem Reading File Contents: ", err)
	}
	fmt.Println("got through to the wait")
	wg.Wait()

}
