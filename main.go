package main

import (
	"MostCommonColors/imageprocessor"
	"bufio"
	"fmt"
	"image/color"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	//open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	//create a map of already seen URLs
	processedURLs := make(map[string]bool)
	urlCount := 0

	//read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//parse the URL
		url := scanner.Text()
		urlCount++
		fmt.Println(" Begin URL #", urlCount, ": ", url)

		//have we seen this one before? if so, skip it
		if _, ok := processedURLs[url]; ok {
			continue
		}
		processedURLs[url] = true
		fmt.Println("Processed URLs: ", len(processedURLs))

		resultChan := make(chan []color.Color)
		errorChan := make(chan error)

		// Launch a goroutine to download and process the image
		wg.Add(1)
		go func() {
			fmt.Println("Downloading and processing ", url)
			imageprocessor.DownloadAndProcessImage(url, resultChan, errorChan)
			wg.Done()
		}()

		// Launch a goroutine to write the results to the CSV file
		wg.Add(1)
		go func() {
			fmt.Println("Writing Results ", url)
			imageprocessor.WriteResultsToCSV(url, <-resultChan, errorChan)
			wg.Done()
		}()
	}
	wg.Wait()

	//did we have a problem scanning?
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file: ", err)
	}
}
