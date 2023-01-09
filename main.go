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
	file, err := os.Open("hugeinput.txt")
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	//read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//parse the URL
		url := scanner.Text()
		fmt.Println("URL: ", url)
		resultChan := make(chan []color.Color)
		errorChan := make(chan error)

		//process the image
		wg.Add(1)

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

	//did we have a problem scanning?
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file: ", err)
	}
	wg.Wait()
}
