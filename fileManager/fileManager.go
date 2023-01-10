package fileManager

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func DownloadAndDecodeImage(url string) (*image.RGBA, error) {
	//download the image
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("problem downloading the image at %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	//decode the image
	jpegImage, err := jpeg.Decode(resp.Body)
	if err != nil {
		fmt.Println("Problem jpeg decoding the file: ", err)
		return nil, err
	}
	// Convert jpeg.Decode image to image.RGBA

	bounds := jpegImage.Bounds()
	rgbaImage := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgbaImage.Set(x, y, jpegImage.At(x, y))
		}
	}

	return rgbaImage, nil
}

func ImportURLs(filename string) (map[string]bool, error) {
	//open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}
	defer file.Close()

	//create a map of already seen URLs
	processedURLs := make(map[string]bool)

	//read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//parse the URL
		inputURL := scanner.Text()

		//clean up the parsed string, does it conform to our expectations?
		u, err := url.Parse(inputURL)
		if err != nil {
			// The input string is not a valid URL
			fmt.Println("The input string is not a valid URL:", err)
			continue
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			// The URL has an unsupported scheme
			fmt.Println("The URL has an unsupported scheme:", u.Scheme)
			continue
		}

		if !strings.HasSuffix(u.Path, ".jpeg") && !strings.HasSuffix(u.Path, ".jpg") {
			// The URL does not end in .jpeg or .jpg
			fmt.Println("The URL is not for a supported image type")
			continue
		}

		// The URL is a valid HTTP URL that ends in .jpeg
		fmt.Println("The URL is a valid HTTP URL that ends in .jpeg")

		//have we seen this one before? if so, skip it
		if _, ok := processedURLs[inputURL]; ok {
			continue
		}
		processedURLs[inputURL] = true
		fmt.Println("Processed URLs: ", len(processedURLs))
	}
	return processedURLs, nil

}
func WriteResultToCSV(fileName string, result string) (bool, error) {
	// Open the CSV file
	file, err := os.OpenFile("output.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("can't open output.csv: ", err)
		return false, err
	}
	defer file.Close()

	// Write the result to the file
	_, err = file.WriteString(result)
	if err != nil {
		fmt.Println("problem writing out to the file: ", err)
		return false, err
	}
	return true, nil
}
